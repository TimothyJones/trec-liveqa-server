package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	//  "log"
	"net/http"
)

// Summarizer implements a simplistic interface for query biased summarizers.
// It takes a Document, a Question, and an integer that represents the
// character limit of the produced summary (this is indicative somehow)
type Summarizer interface {
	GetSummary(docs []Document, q *Question, limit int) string
}

// DummySummarizer returns the full document text as the summary
type DummySummarizer struct{}

func NewDummySummarizer() Summarizer {
	return &DummySummarizer{}
}

func (ls *DummySummarizer) GetSummary(docs []Document, q *Question, limit int) string {
	return docs[0].Text
}

// RemoteSummarizer talks to a remote summarization server in JSON
type RemoteSummarizer struct{ Url string }

func NewRemoteSummarizer(url string) Summarizer {
	return &RemoteSummarizer{Url: url}
}

// RemoteSummarizerRequest simply wraps around the Summarizer protocol
type RemoteSummarizerRequest struct {
	Docs  []Document `json:"documents"`
	Qu    Question   `json:"question"`
	Limit int        `json:"limit"`
}

type RemoteSummarizerResponse struct {
	Summary string `json:"summary"`
}

func (ss *RemoteSummarizer) GetPayload(docs []Document, q *Question, limit int) []byte {
	req := &RemoteSummarizerRequest{docs, *q, limit}
	payload, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return payload
}

func (ss *RemoteSummarizer) GetSummary(docs []Document, q *Question, limit int) string {
	payload := ss.GetPayload(docs, q, limit)
	resp, err := http.Post(ss.Url, "application/json", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var rsr RemoteSummarizerResponse
	if err = json.Unmarshal(body, &rsr); err != nil {
		panic(err)
	}
	return rsr.Summary
}
