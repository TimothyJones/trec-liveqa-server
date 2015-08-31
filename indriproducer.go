package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func Sanitize(r rune) rune {
	switch {
	case unicode.IsPunct(r):
		return ' '
	case unicode.IsMark(r):
		return ' '
	case unicode.IsSymbol(r):
		return ' '
	}
	return r
}

func Truncate(s string, limit int) string {
	var buf bytes.Buffer
	buf.WriteString(s)

	if buf.Len() > limit-3 {
		buf.Truncate(limit - 3)
		return buf.String() + "..."
	} else {
		return buf.String()
	}
}

type IndriAnswerProducer struct {
	Repository     string   `json:"repository"`
	SummarizerUrl  string   `json:"summarizer-url"`
	RunQueryArgs   []string `json:"run-query-args"`
	ExpansionType  string   `json:"expansion-type"`
	ExpansionCount int      `json:"expansion-count"`
}

func NewIndriAnswerProducer(config string) (AnswerProducer, error) {
	ap := &IndriAnswerProducer{}

	byt, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byt, ap); err != nil {
		return nil, err
	}

	if _, err := os.Stat(ap.Repository); err != nil {
		return nil, err
	}

	log.Printf("indriproducer: Repository `%s`\n", ap.Repository)
	log.Printf("indriproducer: SummarizerUrl `%s`\n", ap.SummarizerUrl)
	log.Printf("indriproducer: RunQueryArgs `%s`\n", ap.RunQueryArgs)
	log.Printf("indriproducer: ExpansionType `%s`\n", ap.ExpansionType)
	log.Printf("indriproducer: ExpansionCount `%v`\n", ap.ExpansionCount)
	return ap, nil
}

type IndriQueryResult struct {
	Score       float64
	Docno       string
	StartOffset int
	EndOffset   int
}

// IndriRunQuery executes the query and returns top k docnos
func IndriRunQuery(repo string, query string, k int, args []string) ([]IndriQueryResult, error) {
	query = strings.Map(Sanitize, query)

	var results []IndriQueryResult
	callArgs := append(
		[]string{"-index=" + repo, "-count=" + strconv.Itoa(k), "-query.text=" + query},
		args...)

	out, err := exec.Command("IndriRunQuery", callArgs...).Output()
	if err != nil {
		return results, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		score, _ := strconv.ParseFloat(fields[0], 64)
		docno := fields[1]
		start, _ := strconv.Atoi(fields[2])
		end, _ := strconv.Atoi(fields[3])
		results = append(results, IndriQueryResult{score, docno, start, end})
	}
	return results, nil
}

// RemoveDuplicateDocnos returns a list of unique docnos
func RemoveDuplicateDocnos(docnos []string) []string {
	w := 0

loop:
	for _, docno := range docnos {
		for j := 0; j < w; j++ {
			if docno == docnos[j] {
				continue loop
			}
		}
		docnos[w] = docno
		w++
	}
	return docnos[:w]
}

// IndriDumpText retrieves texts stored in the index
func IndriDumpText(repo string, docno string) (string, error) {
	out, err := exec.Command(
		"dumpindex", repo, "documentid", "docno", docno).Output()
	if err != nil {
		return "", err
	}

	internalDocno := strings.TrimSpace(string(out))

	out, err = exec.Command(
		"dumpindex", repo, "documenttext", internalDocno).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ParseTRECDocument parse texts into Documents
func ParseTRECDocument(text string) (doc Document) {
	//  matchedTags := regexp.MustCompile("</?\\S+(?:\\s+\\S+=\".*?\")*>")
	matchedTags := regexp.MustCompile("</?.*?>")

	lines := strings.Split(strings.TrimSpace(text), "\n")
	var buf bytes.Buffer
	var docno string
	var ok = false
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "<DOCNO>") &&
			strings.HasSuffix(line, "</DOCNO>"):
			docno = strings.TrimSuffix(
				strings.TrimPrefix(line, "<DOCNO>"), "</DOCNO>")
		case strings.HasPrefix(line, "<TEXT>"):
			ok = true
		case strings.HasPrefix(line, "</TEXT>"):
			ok = false
		case ok:
			newline := matchedTags.ReplaceAllString(line, "")
			if len(newline) > 0 {
				buf.WriteString(newline + " ")
			}
		}
	}
	return Document{Docno: docno, Text: buf.String()}
}

func GetPassage(text string, start int, end int) string {
	words := strings.Fields(text)
	if start < 0 {
		start = 0
	}
	if end > len(words) {
		end = len(words)
	}
	return strings.Join(words[start:end], " ")
}

func GetQueryTerms(text string) []string {
	return strings.Fields(strings.Map(Sanitize, strings.ToLower(text)))
}

func PrepareOrdinaryQuery(terms []string) string {
	return strings.Join(terms, " ")
}

func PreparePassageQuery(terms []string) string {
	return fmt.Sprintf(
		"#combine[passage100:50]( %s )", strings.Join(terms, " "))
}

func PrepareSDQuery(terms []string) string {
	var od, ud []string
	for i := 1; i < len(terms); i++ {
		//  if stopwords[terms[i-1]] || stopwords[terms[i]] {
		//  continue
		//  }

		od = append(od, fmt.Sprintf("#1( %s )", strings.Join(terms[i-1:i+1], " ")))
		ud = append(ud, fmt.Sprintf("#uw8( %s )", strings.Join(terms[i-1:i+1], " ")))
	}
	query := fmt.Sprintf(
		"#weight( %1.2f #combine( %s ) %1.2f #combine( %s ) %1.2f #combine ( %s ) )",
		0.85, strings.Join(terms, " "),
		0.10, strings.Join(od, " "),
		0.05, strings.Join(ud, " "),
	)
	return query
}

func (ap *IndriAnswerProducer) GetAnswer(result chan *Answer, q *Question) {
	var answer *Answer
	var summary string
	var resources []string
	var passages []string

	summarizers := []Summarizer{
		NewRemoteSummarizer(ap.SummarizerUrl),
		NewDummySummarizer(),
	}

	timeout := time.After(15 * time.Second)
	headwordchan := GetHeadWord(q.Title)
	expansion := ""
	cache := make(map[string]Document)

HeadWordLoop:
	select {
	case <-timeout:
		log.Println("Query '%s' timed out wating for headword")
		break HeadWordLoop
	case headwords := <-headwordchan:
		switch ap.ExpansionType {
		case "word2vec":
			expansion = word2vec(headwords, ap.ExpansionCount)
		case "wordnet":
			expansion = wordnet(headwords, ap.ExpansionCount)
		}
		log.Printf("Query '%s' has headword(s) '%s'; with synonyms '%s'\n", q.Title, headwords, expansion)
	}

	//  expandedQuery := strings.Join(GetQueryTerms(q.Title+" "+expansion), " ")
	//  queryString := TrimQuery(expandedQuery)
	//  terms := GetQueryTerms(queryString)
	//  query := PreparePassageQuery(terms)
	originalQuery := strings.Join(GetQueryTerms(q.Title), " ")
	queryString := TrimQuery(originalQuery)
	terms := GetQueryTerms(queryString + " " + expansion)
	query := PreparePassageQuery(terms)

	q.Body += " " + expansion

	results, err := IndriRunQuery(ap.Repository, query, 3, ap.RunQueryArgs)
	if err != nil {
		answer = NewErrorAnswer(q, err)
		goto end
	}

	for _, result := range results {
		doc, ok := cache[result.Docno]
		if !ok {
			text, err := IndriDumpText(ap.Repository, result.Docno)
			if err != nil {
				answer = NewErrorAnswer(q, err)
				goto end
			}

			doc = ParseTRECDocument(text)
			cache[result.Docno] = doc
		}

		passages = append(passages,
			GetPassage(doc.Text, result.StartOffset, result.EndOffset))
	}

	for docno := range cache {
		resources = append(resources, docno)
	}

	for _, summarizer := range summarizers {
		summary, err = summarizer.GetSummary(passages, q, config.AnswerSize)
		if err != nil {
			answer = NewErrorAnswer(q, err)
			continue
		}

		answer = &Answer{
			Answered:  "yes",
			Pid:       config.Pid,
			Qid:       q.Qid,
			Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
			Content:   Truncate(summary, config.AnswerSize),
			Resources: strings.Join(resources, ","),
		}
		goto end
	}

end:
	result <- answer
}
