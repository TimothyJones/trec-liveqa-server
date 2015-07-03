package main

import (
	"fmt"
	"time"
)

type AnswerProducer interface {
	GetAnswer(chan *Answer, *Question)
}

type DummyAnswerProducer struct{}

func (*DummyAnswerProducer) GetAnswer(result chan *Answer, q *Question) {
	result <- &Answer{
		Answered:  "yes",
		Pid:       "demo-id-01",
		Qid:       q.Qid,
		Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
		Content:   fmt.Sprintf("Title: %s;  Body: %s;  Category: %s", q.Title, q.Body, q.Category),
		Resources: "resource1,resource2",
	}
}
