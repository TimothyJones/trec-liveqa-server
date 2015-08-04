package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type LogWatch struct{}

func NewLogWatch() (*LogWatch, error) {
	// Set up logging
	logfile, err := os.OpenFile(config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stderr, logfile))
	}

	return &LogWatch{}, err
}

func (lw *LogWatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("tail", "-n", "1000", config.LogPath).Output()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(out)
}
