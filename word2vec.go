package main

import (
	"github.com/TimothyJones/Go-HTTP-JSON-RPC/httpjsonrpc"
	"log"
	"strings"
)

func word2vec(word string) string {
	args := make([]interface{}, 1)
	args[0] = word

	r, err := httpjsonrpc.Call(config.Word2VecServer, "expand", 0, args)
	if err != nil {
		log.Printf("verbose error info: %#v", err)
	}

	if _, ok := r["result"]; !ok {
		log.Println("Word2Vec: not getting any result")
		return ""
	}

	x, ok := r["result"].([]interface{})
	if !ok {
		log.Println("Word2Vec: cannot parse the result")
		return ""
	} else {
		strs := make([]string, len(x))
		for i, item := range x {
			strs[i] = item.(string)
		}
		result := strings.Join(strs, " ")
		log.Printf("Word2Vec expanded '%s' to '%s'\n", word, result)
		return result
	}
}
