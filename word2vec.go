package main

import (
	"github.com/TimothyJones/Go-HTTP-JSON-RPC/httpjsonrpc"
	"log"
	"strings"
)

func word2vec(str string) string {
	args := make([]interface{}, 1)
	args[0] = str

	r, err := httpjsonrpc.Call(config.Word2VecServer, "expand", 0, args)
	if err != nil {
		log.Printf("Unable to get word2vec expansion for '%s'\n", str)
		return ""
	}

	if _, ok := r["result"]; !ok {
		return ""
	}

	x, ok := r["result"].([]string)
	if ok {
		return strings.Join(x, " ")
	} else {
		return ""
	}

	return str
}
