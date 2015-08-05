package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"net/http"
	"net/url"
	"strings"
	"time"
)

// GalagoAnswerProducer forwards the question on to Galago
// It implements the AnswerProducer interface
type GalagoAnswerProducer struct {
	Host string `json:"host"`
}

func NewGalagoAnswerProducer(filename string) (AnswerProducer, error) {
	ap := &GalagoAnswerProducer{}

	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byt, ap); err != nil {
		return nil, err
	}

	return ap, nil
}

func (p *GalagoAnswerProducer) GetUrl(query string) (*url.URL, error) {
	u, err := url.Parse(p.Host)
	if err != nil {
		return u, err
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"

	}
	u.Path = u.Path + "searchxml"

	if u.Scheme == "" {
		u.Scheme = "http"
	}
	uq := u.Query()
	uq.Set("q", query)
	uq.Set("n", "20")
	u.RawQuery = uq.Encode()

	return u, nil
}

// GetAnswer submits the question to a Galago index and returns an answer
func (p *GalagoAnswerProducer) GetAnswer(result chan *Answer, q *Question) {
	var answer *Answer

	url, err := p.GetUrl(q.Title)
	if err != nil {
		answer = NewErrorAnswer(q, err)
		goto end
	}

	// TODO: Read the content out of the response

	result <- &Answer{
		Answered:  "yes",
		Pid:       config.Pid,
		Qid:       q.Qid,
		Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
		Content:   fmt.Sprintf("%s Title: %s;  Body: %s;  Category: %s", url.String(), q.Title, q.Body, q.Category),
		Resources: "resource1,resource2",
	}

end:
	result <- answer
}
