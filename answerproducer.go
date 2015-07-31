package main

import (
	"fmt"
	"time"
)

// The AnswerProducer inteface allows objects to answer questions in a way that
// the LiveQA handler can understand
type AnswerProducer interface {
	GetAnswer(chan *Answer, *Question)
}

// DummyAnswerProducer is a place holder for instantly returning answers that just parrot the question
type DummyAnswerProducer struct{}

// GetAnswer returns the question title, body and category as the answer.
func (*DummyAnswerProducer) GetAnswer(result chan *Answer, q *Question) {
	result <- &Answer{
		Answered:  "yes",
		Pid:       config.Pid,
		Qid:       q.Qid,
		Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
		Content:   fmt.Sprintf("Title: %s;  Body: %s;  Category: %s", q.Title, q.Body, q.Category),
		Resources: "resource1,resource2",
	}
}
