package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

type Question struct {
	Qid          string `schema:"qid" json:"qid"`
	Title        string `schema:"title" json:"title"`
	Body         string `schema:"body" json:"body"`
	Category     string `schema:"category" json:"category"`
	ReceivedTime time.Time
}

func (q *Question) String() string {
	b, err := json.Marshal(q)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(b)
}

type AnswerWrapper struct {
	XMLName xml.Name `xml:"xml"`
	Answer  *Answer  `xml:"answer"`
}

type Answer struct {
	XMLName   xml.Name `xml:"answer"`
	Answered  string   `xml:"answered,attr"`
	Pid       string   `xml:"pid,attr"`
	Qid       string   `xml:"qid,attr"`
	Time      int64    `xml:"time,attr"`
	Content   string   `xml:"content"`
	Resources string   `xml:"resources"`
}

func (a *AnswerWrapper) String() string {
	output, err := xml.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(output)
}

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
