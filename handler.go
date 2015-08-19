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

// The LiveQA struct provides a handler that processes Questions submitted via
// GET or POST. It implements http.Handle
type LiveQA struct {
	Producers []AnswerProducer
}

// NewLiveQA creates a LiveQA structure
func NewLiveQA() *LiveQA {
	return &LiveQA{
		Producers: make([]AnswerProducer, 0, 10),
	}
}

// AddProcessor adds an AnswerProducer to this handler. All AnswerProcessors
// are queried when a question is processed, and the most recent answer (if
// any) is returned by ProcessQuestion
func (lqa *LiveQA) AddProducer(ap AnswerProducer) {
	lqa.Producers = append(lqa.Producers, ap)
}

// ProcessQuery processes a question and returns a wrapped answer.
// It submits the question to all AnswerProducers in the Producers slice, and
// returns the most recent answer. If none of the AnswerProducers return in the
// timeout, it returns a failed answer that contains a little information about
// the question. If there are still answers being produced at the time the
// timeout hits, the most recent answer is returned. If all AnswerProducers
// return before the timeout, the function will return early.
func (lqa *LiveQA) ProcessQuestion(q *Question) *AnswerWrapper {
	answers := make(chan *Answer, 1)

	// Kick off all the answer producers
	var wg sync.WaitGroup
	for _, ap := range lqa.Producers {
		go ap.GetAnswer(answers, q)
		wg.Add(1)
	}

	// We want to be able to exit after a timeout
	timeout := time.After(time.Duration(config.Timeout) * time.Second)
	answer := NewTimeoutAnswer(q, config.Timeout)

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

func (lqa *LiveQA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Drop query if without a proper Qid
	if len(q.Qid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("IP", r.RemoteAddr, "QID", q.Qid)

	// Process query here
	a := lqa.ProcessQuestion(q)

	log.Println("Got query `", q.Title, "` for", q.Qid)
	log.Println("Got answer `", a.Answer.Content, "` with resources `",
		a.Answer.Resources, "` for", q.Qid, "in time", a.Answer.Time)

	fmt.Fprintf(w, "%s%s\n", xml.Header, a)
}
