package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
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

	log.Println("QID", q.Qid)

	// Process query here

	a := &AnswerWrapper{
		Answer: &Answer{
			Answered:  "yes",
			Pid:       "demo-id-01",
			Qid:       q.Qid,
			Time:      int64(time.Since(start) / time.Millisecond),
			Content:   fmt.Sprintf("Title: %s;  Body: %s;  Category: %s", q.Title, q.Body, q.Category),
			Resources: "resource1,resource2",
		},
	}

	log.Println("Got answer `", a.Answer.Content, "` for", q.Qid, "in time", a.Answer.Time)

	fmt.Fprintf(w, "%s%s\n", xml.Header, a)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
