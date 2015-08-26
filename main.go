package main

import (
	"github.com/vharitonsky/iniflags"
	"log"
	"net/http"
	"strconv"
)

var factory map[string]func(string) (AnswerProducer, error)

func init() {
	factory = make(map[string]func(string) (AnswerProducer, error))
	factory["indri"] = NewIndriAnswerProducer
	factory["dummy"] = NewDummyAnswerProducer
	factory["galago"] = NewGalagoAnswerProducer
}

func main() {
	iniflags.Parse()

	// Create a liveQA
	lqa := NewLiveQA()

	// Add answer producers
	count := 0
	for _, name := range config.Producers {
		if f, ok := factory[name]; ok {
			ap, err := f(name + ".json")
			if err != nil {
				log.Printf("[Error initialising %s] %s\n", name, err)
			} else {
				lqa.AddProducer(ap)
				log.Printf("Initialised '%s' answer producer\n", name)
			}
		} else {
			log.Printf("[Error initialising %s] Unrecognised answer producer '%s'\n", name, name)
		}
	}
	log.Printf("Initialised %d of %d total answer producers\n", count, len(config.Producers))

	http.Handle("/", lqa)

	// Set up logging
	if lw, err := NewLogWatch(); err != nil {
		log.Printf("[flag -logfile] '%s' will carry on without the logfile\n", err)
	} else {
		http.Handle("/tail1000", lw)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
