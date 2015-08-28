package main

import (
	"flag"
	"strings"
)

var config struct {
	Pid            string
	Timeout        int
	Port           int
	LogPath        string
	Producers      producers
	AnswerSize     int
	HeadWordServer string
	Word2VecServer string
	KrovetzBinary  string
	MaxScores      string
}

type producers []string

func (i *producers) Set(value string) error {
	for _, proc := range strings.Split(value, ",") {
		*i = append(*i, proc)
	}
	return nil
}

func (i *producers) String() string {
	return strings.Join(*i, ",")
}

func init() {
	config.Producers = make([]string, 0, 0)
	flag.IntVar(&config.Timeout, "timeout", 55, "timeout, in seconds")
	flag.IntVar(&config.Port, "port", 8080, "HTTP service port")
	flag.StringVar(&config.Pid, "pid", "demo-pid-01", "participant ID")
	flag.StringVar(&config.HeadWordServer, "headwordserver", "http://localhost:8001", "Location of headword server")
	flag.StringVar(&config.Word2VecServer, "word2vecserver", "http://localhost:8002", "Location of word2vec server")
	flag.StringVar(&config.LogPath, "log", "liveqa.log", "path to log file")
	flag.StringVar(&config.MaxScores, "maxscoresfile", "/collection/jsc/for_liveqa/ordered_max_scores.txt", "path to max scores file")
	flag.StringVar(&config.KrovetzBinary, "krovetzbin", "/storage/rc/shane-reduce/krovetz-query", "path to krovetz stemmer binary")
	flag.Var(&config.Producers, "producer", "comma separated list of processors to use on this server")
	flag.IntVar(&config.AnswerSize, "answersize", 1000, "answer size limit, in number of chars")
}
