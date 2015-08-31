package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Summarizer implements a simplistic interface for query biased summarizers.
// It takes a list of passages, a Question, and an integer that represents the
// character limit of the produced summary (this is indicative somehow)
type Summarizer interface {
	GetSummary(passages []string, q *Question, limit int) (string, error)
}

// DummySummarizer returns the full document text as the summary
type DummySummarizer struct{}

func NewDummySummarizer() Summarizer {
	return &DummySummarizer{}
}

func (ls *DummySummarizer) GetSummary(passages []string, q *Question, limit int) (string, error) {
	return passages[0], nil
}

// RemoteSummarizer talks to a remote summarization server in JSON
type RemoteSummarizer struct{ Url string }

func NewRemoteSummarizer(url string) Summarizer {
	return &RemoteSummarizer{Url: url}
}

// RemoteSummarizerRequest simply wraps around the Summarizer protocol
type RemoteSummarizerRequest struct {
	Texts    []string `json:"texts"`
	Question Question `json:"question"`
	Limit    int      `json:"limit"`
}

type RemoteSummarizerResponse struct {
	Summary string `json:"summary"`
}

func (ss *RemoteSummarizer) GetPayload(passages []string, q *Question, limit int) []byte {
	req := &RemoteSummarizerRequest{passages, *q, limit}
	payload, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return payload
}

func (ss *RemoteSummarizer) GetSummary(passages []string, q *Question, limit int) (string, error) {
	var summary string
	payload := ss.GetPayload(passages, q, limit)
	resp, err := http.Post(ss.Url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return summary, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return summary, err
	}

	var rsr RemoteSummarizerResponse
	if err = json.Unmarshal(body, &rsr); err != nil {
		return summary, err
	}
	return rsr.Summary, nil
}
