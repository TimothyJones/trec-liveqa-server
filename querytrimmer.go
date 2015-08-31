package main

import (
	"bufio"
	//"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var importance map[string]float64

func StemQuery(query string) string {
	cmd := exec.Command(config.KrovetzBinary)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	io.WriteString(stdin, query)
	if err := stdin.Close(); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	output := string(bytes)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(output)
}

type QueryTerm struct {
	Term       string
	Importance float64
}

func (q QueryTerm) String() string {
	return q.Term
}

type ByImportance []QueryTerm

func (a ByImportance) Len() int           { return len(a) }
func (a ByImportance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByImportance) Less(i, j int) bool { return a[i].Importance > a[j].Importance }

func TrimQuery(query string) string {
	if !config.TrimQueries {
		return query
	}

	stemmedQuery := query //StemQuery(query)
	stemmedTerms := strings.Split(stemmedQuery, " ")
	terms := strings.Split(query, " ")
	if len(terms) != len(stemmedTerms) {
		log.Printf("ERROR '%s' does not have the same number of terms after stemming to '%s'\n", query, stemmedQuery)
		return query
	}
	qts := make([]QueryTerm, len(terms), len(terms))

	for i, _ := range terms {
		if _, ok := stopwords[terms[i]]; ok {
			// This term is stopped
			continue
		}
		imp, ok := importance[stemmedTerms[i]]
		if !ok {
			log.Printf("WARN: '%s' is not a term in the imporance map\n", stemmedTerms[i])
			imp = 0
		}
		qt := QueryTerm{Term: terms[i], Importance: imp}
		qts = append(qts, qt)
	}

	sort.Sort(ByImportance(qts))
	/*
		res := "Terms: "
		for _, term := range qts {
			res += term.Term + "(" + fmt.Sprintf("%f", term.Importance) + ") "
		}
		log.Println(res)              */

	trimmed := qts
	if len(trimmed) > config.WordsPerQuery {
		trimmed = trimmed[:config.WordsPerQuery]
	}

	result := ""
	for _, term := range trimmed {
		result += term.Term + " "
	}
	log.Printf("Query trimmed down to `%s`\n", result)
	return result
}

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
