package main

import (
	"net/http"
	"os/exec"
)

type LogWatch struct {
	Filename string
}

func NewLogWatch(filename string) *LogWatch {
	return &LogWatch{Filename: filename}
}

func (lw *LogWatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("tail", "-n", "1000", lw.Filename).Output()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(out)
}
