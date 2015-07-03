package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"sync"
	"time"
)

type LiveQA struct {
	Producers []AnswerProducer
	timeout   int
}

// NewLiveQA creates a LiveQA structure, with the specified timeout
func NewLiveQA(timeout int) *LiveQA {
	return &LiveQA{
		Producers: make([]AnswerProducer, 0, 10),
		timeout:   timeout,
	}
}

// AddProcessor adds a processor to this handler
func (lqa *LiveQA) AddProcessor(ap AnswerProducer) {
	lqa.Producers = append(lqa.Producers, ap)
}

// ProcessQuery processes a question and returns a wrapped answer
func (lqa *LiveQA) ProcessQuery(q *Question) *AnswerWrapper {
	answers := make(chan *Answer, 1)

	// Kick off all the answer producers
	var wg sync.WaitGroup
	for _, ap := range lqa.Producers {
		go ap.GetAnswer(answers, q)
		wg.Add(1)
	}

	// We want to be able to exit after a timeout
	timeout := time.After(time.Duration(lqa.timeout) * time.Second)
	answer := NewTimeOutAnswer(q, lqa.timeout)

	// We also want to be able to exit after all producers have returned
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Get the most recent answer, or timeout
Loop:
	for {
		select {
		case answer = <-answers:
			wg.Done()
		case <-timeout:
			break Loop
		case <-done:
			break Loop
		}
	}

	a := &AnswerWrapper{Answer: answer}
	return a
}

func (lqa *LiveQA) handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	q := &Question{}
	decoder := schema.NewDecoder()
	err = decoder.Decode(q, r.Form)
	if err != nil {
		log.Println(err)
		return
	}
	q.ReceivedTime = time.Now()

	log.Println("QID", q.Qid)

	// Process query here
	a := lqa.ProcessQuery(q)

	log.Println("Got answer `", a.Answer.Content, "` for", q.Qid, "in time", a.Answer.Time)

	fmt.Fprintf(w, "%s%s\n", xml.Header, a)
}
