package main

import (
	"flag"
	"strings"
)

var config struct {
	Pid        string
	Timeout    int
	Port       int
	LogPath    string
	Processors processors
}

type processors []string

func (i *processors) Set(value string) error {
	for _, proc := range strings.Split(value, ",") {
		*i = append(*i, proc)
	}
	return nil
}

func (i *processors) String() string {
	return strings.Join(*i, ",")
}

func init() {
	flag.IntVar(&config.Timeout, "timeout", 30, "timeout, in seconds")
	flag.IntVar(&config.Port, "port", 8080, "HTTP service port")
	flag.StringVar(&config.Pid, "pid", "demo-pid-01", "participant ID")
	flag.StringVar(&config.LogPath, "log", "liveqa.log", "path to log file")
	flag.Var(&config.Processors, "processor", "comma separated list of processors to use on this server")
}
