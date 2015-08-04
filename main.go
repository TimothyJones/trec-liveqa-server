package main

import (
	"flag"
	"github.com/vharitonsky/iniflags"
	"log"
	"net/http"
	"strconv"
)

func main() {
	index := flag.String("index", "", "path to Indri index")
	iniflags.Parse()

	// Create a liveQA
	lqa := NewLiveQA()

	// Add a dummy answer producer to it
	lqa.AddProducer(&DummyAnswerProducer{})

	// Add an indri index answer producer
	iap, err := NewIndriIndexAnswerProducer(*index)
	if err != nil {
		log.Fatal("[indri index]", err)
	}
	lqa.AddProducer(iap)

	http.Handle("/", lqa)

	// Set up logging
	if lw, err := NewLogWatch(); err != nil {
		log.Printf("[flag -logfile] '%s' will carry on without the logfile\n", err)
	} else {
		http.Handle("/tail1000", lw)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
