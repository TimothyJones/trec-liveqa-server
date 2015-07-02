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

	log.Println("Query", q.Qid)

	// Process query here

	a := &AnswerWrapper{
		Answer: &Answer{
			Answered:  "yes",
			Pid:       "demo-id-01",
			Qid:       q.Qid,
			Time:      int64(time.Since(start) / time.Millisecond),
			Content:   "[YOUR ANSWER HERE]",
			Resources: "resource1,resource2",
		},
	}

	log.Println("Query", q.Qid, "got answer `", a.Answer.Content, "` in time", a.Answer.Time)

	fmt.Fprintf(w, "%s%s", xml.Header, a)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
