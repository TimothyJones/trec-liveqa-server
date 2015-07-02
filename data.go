package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
)

type Question struct {
	Qid      string `schema:"qid" json:"qid"`
	Title    string `schema:"title" json:"title"`
	Body     string `schema:"body" json:"body"`
	Category string `schema:"category" json:"category"`
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
	output, err := xml.MarshalIndent(a, "  ", "    ")
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return string(output)
}
