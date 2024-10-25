package query

import (
	"iter"
	"maps"
	"strings"
)

type Query struct {
	MetricName string
	Labels     map[string]string
}

func Parse(s string) (Query, error) {
	if idx := strings.IndexRune(s, '{'); idx > -1 {
		var q Query
		q.MetricName = s[0:idx]
		q.Labels = maps.Collect(labelpairs(s[idx:]))
		return q, nil
	}

	return Query{MetricName: s}, nil
}

func labelpairs(s string) iter.Seq2[string, string] {
	pairs := strings.Split(s[1:len(s)-1], ",")
	return func(yield func(string, string) bool) {
		for _, p := range pairs {
			key, val, found := strings.Cut(p, "=")
			if !yield(strings.TrimSpace(key), strings.TrimSpace(val)) || !found {
				break
			}
		}
	}
}
