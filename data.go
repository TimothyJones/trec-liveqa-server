package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

// Question contains all the information from a question requested via HTTP,
// plus a field (filled by the reciever) for the time the request was received.
type Question struct {
	Qid          string    `schema:"qid" json:"qid"`           // The query ID
	Title        string    `schema:"title" json:"title"`       // The query text
	Body         string    `schema:"body" json:"body"`         // The body of the question (may be empty)
	Category     string    `schema:"category" json:"category"` // The category of the question (may be empty)
	ReceivedTime time.Time // The time this quesiton was recieved
}

// AsJSON returns a json representation of the Question structure
func (q *Question) AsJSON() string {
	b, err := json.Marshal(q)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(b)
}

// String is a convieniencec returns the json representation of the Question structure
func (q *Question) String() string {
	return q.AsJSON()
}

// AnswerWrapper provides the xml wrapper for an answer. It's only needed at
// serilisation time
type AnswerWrapper struct {
	XMLName xml.Name `xml:"xml"`
	Answer  *Answer  `xml:"answer"`
}

// Answer provides the representation of an answer from a AnswerProcessor. It
// can be serilised to xml using xml.Marshal(), or the convienience method answer.AsXML()
type Answer struct {
	XMLName   xml.Name `xml:"answer"`
	Answered  string   `xml:"answered,attr"`
	Pid       string   `xml:"pid,attr"`
	Qid       string   `xml:"qid,attr"`
	Time      int64    `xml:"time,attr"`
	Content   string   `xml:"content"`
	Resources string   `xml:"resources"`
}

// AsXML returns a XML representation of the Answer structure, suitable for
// returning to the TREC question client
func (a *AnswerWrapper) AsXML() string {
	output, err := xml.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(output)
}

func (a *AnswerWrapper) String() string {
	return a.AsXML()
}

// NewTimeOutAnswer provides the error response "answer" when a timeout has
// been hit. The timeout supplied is included in the response
func NewTimeOutAnswer(q *Question, timeout int) *Answer {
	return &Answer{
		Answered:  "no",
		Pid:       "demo-id-01",
		Qid:       q.Qid,
		Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
		Content:   fmt.Sprintf("TIMEOUT after %d seconds: Title: %s;  Body: %s;  Category: %s", timeout, q.Title, q.Body, q.Category),
		Resources: "resource1,resource2",
	}
}
