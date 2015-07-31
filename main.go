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

	// Create a liveQA listener with a timeout of 30 seconds
	lqa := NewLiveQA()

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
	logfile, err := os.OpenFile(config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stderr, logfile))

		lw := NewLogWatch(config.LogPath)
		http.Handle("/tail1000", lw)
	} else {
		log.Fatalf("[flag -logfile] %s", err)
		log.Fatalln("will carry on without the logfile")
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
