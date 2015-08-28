package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Repository    string `json:"repository"`
	SummarizerUrl string `json:"summarizer-url"`
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
	return ap, nil
}

// IndriRunQuery executes the query and returns top k docnos
func IndriRunQuery(repo string, query string, k int) ([]string, error) {
	query = strings.Map(Sanitize, query)

	var docnos []string
	out, err := exec.Command(
		"IndriRunQuery", "-index="+repo, "-trecFormat=1",
		"-count="+strconv.Itoa(k), "-rule=method:okapi",
		"-fbDocs=10", "-fbTerms=5",
		"-query.text="+query).Output()
	if err != nil {
		return docnos, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		docnos = append(docnos, fields[2])
	}
	return docnos, nil
}

// IndriDumpText retrieves texts stored in the index
func IndriDumpText(repo string, docnos []string) ([]string, error) {
	var texts []string
	for _, docno := range docnos {
		out, err := exec.Command(
			"dumpindex", repo, "documentid", "docno", docno).Output()
		if err != nil {
			return texts, err
		}
		internalDocno := strings.TrimSpace(string(out))

		out, err = exec.Command(
			"dumpindex", repo, "documenttext", internalDocno).Output()
		if err != nil {
			return texts, err
		}
		texts = append(texts, string(out))
	}
	return texts, nil
}

// ParseTRECDocument parse texts into Documents
func ParseTRECDocument(texts []string) (docs []Document) {
	//  matchedTags := regexp.MustCompile("</?\\S+(?:\\s+\\S+=\".*?\")*>")
	matchedTags := regexp.MustCompile("</?.*?>")

	for _, text := range texts {
		lines := strings.Split(strings.TrimSpace(text), "\n")
		var buf bytes.Buffer
		var docno string
		var ok = false
		for _, line := range lines {
			switch {
			case buf.Len() > 1000:
				break
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
		docs = append(docs, Document{Docno: docno, Text: buf.String()})
	}
	return
}

func IndriGetTopDocument(repo string, query string, k int) ([]Document, error) {
	var docs []Document
	docnos, err := IndriRunQuery(repo, strings.Map(Sanitize, query), k)
	if err != nil {
		return docs, err
	}

	texts, err := IndriDumpText(repo, docnos)
	if err != nil {
		return docs, err
	}
	docs = ParseTRECDocument(texts)
	return docs, nil
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
	var docnos, texts []string
	var docs []Document

	summarizers := []Summarizer{
		NewRemoteSummarizer(ap.SummarizerUrl),
		NewDummySummarizer(),
	}

	timeout := time.After(30 * time.Second)
	headwordchan := GetHeadWord(q.Title)
	expansion := ""

HeadWordLoop:
	select {
	case <-timeout:
		log.Println("Query '%s' timed out wating for headword")
		break HeadWordLoop
	case headwords := <-headwordchan:
		expansion := word2vec(headwords)
		log.Printf("Query '%s' has headword(s) '%s'; with synonyms '%s'\n", q.Title, headwords, expansion)
	}

	queryString := TrimQuery(q.Title + " " + expansion)
	terms := GetQueryTerms(queryString)
	query := PreparePassageQuery(terms)
	//  if len(terms) > 0 {
	//  query = PrepareOrdinaryQuery(terms)
	//  } else {
	//  query = PrepareSDQuery(terms)
	//  }

	docnos, err := IndriRunQuery(ap.Repository, query, 3)
	if err != nil {
		answer = NewErrorAnswer(q, err)
		goto end
	}

	texts, err = IndriDumpText(ap.Repository, docnos)
	if err != nil {
		answer = NewErrorAnswer(q, err)
		goto end
	}

	docs = ParseTRECDocument(texts)
	if len(docs) < 1 {
		answer = NewErrorAnswer(q, errors.New("No answer found"))
		goto end
	}

	for _, doc := range docs {
		resources = append(resources, doc.Docno)
	}

	for _, summarizer := range summarizers {
		summary, err = summarizer.GetSummary(docs, q, config.AnswerSize)
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
