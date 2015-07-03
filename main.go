package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a liveQA listener with a timeout of 3 seconds
	lqa := NewLiveQA(3)
	// Add a dummy answer producer to it
	lqa.AddProducer(&DummyAnswerProducer{})

	http.Handle("/", lqa)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
