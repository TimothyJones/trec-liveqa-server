package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	timeout := flag.Int("timeout", 30, "timeout, in seconds")
	port := flag.Int("port", 8080, "HTTP service port")
	logpath := flag.String("log", "liveqa.log", "path to log file")
	index := flag.String("index", "", "path to Indri index")
	flag.Parse()

	// Create a liveQA listener with a timeout of 30 seconds
	lqa := NewLiveQA(*timeout)

	// Add a dummy answer producer to it
	lqa.AddProducer(&DummyAnswerProducer{})

	// Add an indri index answer producer
	if _, err := os.Stat(*index); err != nil {
		log.Fatalf("[flag -index] %s", err)
		return
	}
	lqa.AddProducer(&IndriIndexAnswerProducer{*index})

	http.Handle("/", lqa)

	// Set up logging
	logfile, err := os.OpenFile(*logpath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stderr, logfile))
	} else {
		log.Fatalf("[flag -logfile] %s", err)
		log.Fatalln("will carry on without the logfile")
	}

	lw := NewLogWatch(*logpath)
	http.Handle("/tail1000", lw)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
