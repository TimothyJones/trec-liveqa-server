package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var importance map[string]float64

func LoadImportance() {
	importance = make(map[string]float64)

	file, err := os.Open(config.MaxScores)
	if err != nil {
		log.Fatalf("Unable to open max scores file '%s': %s\n", config.MaxScores, err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	log.Println("Initilaising Term Importance Map (takes around 1 min with default file)")
	for scanner.Scan() {
		line := scanner.Text()
		A := strings.Split(line, " ")
		if len(A) != 2 {
			log.Printf("MaxScores: Bad line '%s'\n", line)
			continue
		}
		if s, err := strconv.ParseFloat(A[1], 64); err == nil {
			importance[A[0]] = s
		} else {
			log.Printf("MaxScores: Bad float in line '%s'\n", line)
		}
	}
	log.Println("Term Importance Map Complete")
}
