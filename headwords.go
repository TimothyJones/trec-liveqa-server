package main

import (
	"github.com/TimothyJones/Go-HTTP-JSON-RPC/httpjsonrpc"
	"log"
)

func GetHeadWord(question string) chan string {
	c := make(chan string, 1)

	go func() {
		args := make([]interface{}, 1)
		args[0] = []string{question}

		r, err := httpjsonrpc.Call(config.HeadWordServer, "extract_headchunks", 0, args)
		if err != nil {
			log.Printf("Unable to get headword for '%s'\n", question)
			c <- ""
			return
		}

		if _, ok := r["result"]; !ok {
			// We didn't receive a result
			log.Printf("Headword: HeadWord at '%s' question '%s' had no result from the server (%s, %T)\n", config.HeadWordServer, question, r, r)
			c <- ""
			return
		}

		x, ok := r["result"].([]interface{})
		if ok {
			str, ok := x[0].(string)
			if ok {
				c <- str
			} else {
				log.Printf("Headword: Not string (%s, %T)\n", x[0], x[0])
				// We didn't receive a string
				c <- ""
			}
		} else {
			log.Printf("Headword: Not array (%s, %T)\n", r["result"], r["result"])
			// we didn't recieve a result
			c <- ""
		}
	}()
	return c
}
