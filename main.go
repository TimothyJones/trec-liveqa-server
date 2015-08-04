package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var config struct {
	Pid     string
	Timeout int
	Port    int
	LogPath string
}

func main() {
	// Set up config
	flag.IntVar(&config.Timeout, "timeout", 30, "timeout, in seconds")
	flag.IntVar(&config.Port, "port", 8080, "HTTP service port")
	flag.StringVar(&config.Pid, "pid", "demo-pid-01", "participant ID")
	flag.StringVar(&config.LogPath, "log", "liveqa.log", "path to log file")
	index := flag.String("index", "", "path to Indri index")
	flag.Parse()

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
	logfile, err := os.OpenFile(config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stderr, logfile))

		lw := NewLogWatch(config.LogPath)
		http.Handle("/tail1000", lw)
	} else {
		log.Printf("[flag -logfile] '%s' will carry on without the logfile\n", err)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
