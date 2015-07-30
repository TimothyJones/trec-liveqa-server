package main

import (
	"bytes"
	"os/exec"
	"regexp"
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

// Run the query, and dump the top-1 document content
func IndriGetTopDocument(repo string, query string) string {
	query = strings.Map(Sanitize, query)

	out, err := exec.Command(
		"IndriRunQuery", "-index="+repo,
		"-trecFormat=1", "-count=1",
		"-query.text="+query).Output()
	if err != nil {
		return err.Error()
	}

	fields := strings.Fields(string(out))
	if len(fields) < 2 {
		return "[ERROR] No result"
	}

	docno := fields[2]

	out, err = exec.Command(
		"dumpindex", repo, "documentid", "docno", docno).Output()
	if err != nil {
		return err.Error()
	}
	internal_docno := strings.TrimSpace(string(out))

	out, err = exec.Command(
		"dumpindex", repo, "documenttext", internal_docno).Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func Summarize(content string, limit int) string {
	matchedTags := regexp.MustCompile("</?\\w+(?:\\s+\\w+=\".*?\")*>")

	lines := strings.Split(content, "\n")
	var buf bytes.Buffer
	var ok = false
	for _, line := range lines {
		switch {
		case buf.Len() > limit:
			break
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

	if buf.Len() > limit-5 {
		buf.Truncate(limit - 5)
		return buf.String() + "..."
	} else {
		return buf.String()
	}
}

type IndriIndexAnswerProducer struct {
	Repository string
}

func (ap *IndriIndexAnswerProducer) GetAnswer(result chan *Answer, q *Question) {
	content := IndriGetTopDocument(ap.Repository, q.Title)
	summary := Summarize(content, 250)
	result <- &Answer{
		Answered:  "yes",
		Pid:       "demo-id-02",
		Qid:       q.Qid,
		Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
		Content:   summary,
		Resources: "resource1,resource2",
	}
}
