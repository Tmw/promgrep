package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/tmw/promgrep/pkg/exposition"
	"github.com/tmw/promgrep/pkg/query"
)

func main() {
	tokenizer := exposition.NewTokenizer(bufio.NewReader(os.Stdin))
	parser := exposition.NewParser(tokenizer.Tokens())

	entries, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		query, err := query.Parse(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		entries = filter(entries, query)
	}

	for _, ent := range entries {
		fmt.Println(ent.String())
	}
}

func filter(entries []exposition.MetricFamily, q query.Query) []exposition.MetricFamily {
	res := []exposition.MetricFamily{}

	for _, entry := range entries {
		if q.MetricName != "" {
			if q.MetricName != entry.Name {
				continue
			}
		}

		if !labelsMatch(entry, q) {
			continue
		}

		res = append(res, entry)
	}

	return res
}

func labelsMatch(entry exposition.MetricFamily, q query.Query) bool {
	for k, v := range q.Labels {
		if val, ok := entry.Labels[k]; !ok || val != v {
			return false
		}
	}

	return true
}
