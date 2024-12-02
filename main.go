package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmw/promgrep/pkg/exposition"
	"github.com/tmw/promgrep/pkg/metricfamily"
	"github.com/tmw/promgrep/pkg/query"
)

func main() {
	tokenizer := exposition.NewTokenizer(bufio.NewReader(os.Stdin))
	entries, err := metricfamily.Parse(tokenizer.Tokens())

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		allArgs := strings.Join(os.Args[1:], " ")
		q, err := query.Compile(allArgs)
		if err != nil {
			log.Fatal(err)
		}

		entries = filter(entries, q)
	}

	for _, ent := range entries {
		fmt.Println(ent.String())
	}
}

func filter(entries []metricfamily.MetricFamily, q query.Query) []metricfamily.MetricFamily {
	res := []metricfamily.MetricFamily{}

	for _, entry := range entries {
		if !q.MetricName.Match(entry.Name) {
			continue
		}

		if !labelsMatch(entry, q) {
			continue
		}

		res = append(res, entry)
	}

	return res
}

func labelsMatch(entry metricfamily.MetricFamily, q query.Query) bool {
	for k, v := range q.Labels {
		if val, ok := entry.Labels[k]; !ok || !v.Match(val) {
			return false
		}
	}

	return true
}
