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

	// Set up logging
	if lw, err := NewLogWatch(); err != nil {
		log.Printf("[flag -logfile] '%s' will carry on without the logfile\n", err)
	} else {
		http.Handle("/tail1000", lw)
	}

	if config.TrimQueries {
		log.Printf("Testing stemming: '%s'\n", StemQuery("I like to stemming my queries"))
		// This is a function rather an an init, as it relies of the results of the flags
		LoadImportance()
		log.Println(TrimQuery("This is a stemming queries that needs a bit of a trim dream babies"))
	}
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
				count++
			}
		} else {
			log.Printf("[Error initialising %s] Unrecognised answer producer '%s'\n", name, name)
		}
	}
	log.Printf("Initialised %d of %d total answer producers\n", count, len(config.Producers))

	http.Handle("/", lqa)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
