package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Create a liveQA listener with a timeout of 3 seconds
	lqa := NewLiveQA(3)
	// Add a dummy answer producer to it
	lqa.AddProcessor(&DummyAnswerProducer{})

	r.HandleFunc("/", lqa.handler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
