package main

import (
	"testing"
	"time"
)

func TestPayload(t *testing.T) {
	ss := StandardSummarizer{"localhost:8001"}

	doc := &Document{"Trec LiveQA Trac: 2015 The automated question answering (QA) track, which has been one of the most popular tracks in TREC for many years, has focused on the task of providing automatic answers for human questions."}

	q := &Question{
		Qid:          "fake-qid",
		Title:        "What is artificial intelligence?",
		Body:         "history and evolution of artificial intelligence",
		Category:     "Programming & Design",
		ReceivedTime: time.Unix(0, 0),
	}

	expected := `{"document":{"Text":"Trec LiveQA Trac: 2015 The automated question answering (QA) track, which has been one of the most popular tracks in TREC for many years, has focused on the task of providing automatic answers for human questions."},"question":{"qid":"fake-qid","title":"What is artificial intelligence?","body":"history and evolution of artificial intelligence","category":"Programming \u0026 Design","ReceivedTime":"1970-01-01T10:00:00+10:00"},"limit":250}`

	payload := ss.GetPayload(doc, q, 250)
	if string(payload) != expected {
		t.Errorf("Expected payload `%s` but got `%s`", expected, string(payload))
	}
}
