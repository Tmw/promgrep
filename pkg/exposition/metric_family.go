package exposition

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strconv"
	"strings"
)

type MetricFamily struct {
	Name   string
	Val    float64
	Labels map[string]string
}

func sorted[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)

	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				break
			}
		}
	}
}

func (m *MetricFamily) String() string {
	var b strings.Builder
	b.WriteString(m.Name)

	numLabels := len(m.Labels)
	curLabel := 0

	if numLabels > 0 {
		b.WriteRune('{')
		for k, v := range sorted(m.Labels) {
			b.WriteString(fmt.Sprintf("%s=\"%s\"", k, v))

			// not last label? separate with comma and space
			if curLabel != numLabels-1 {
				b.WriteString(", ")
			}

			curLabel++

		}
		b.WriteRune('}')
	}

	b.WriteRune(' ')
	b.WriteString(strconv.FormatFloat(m.Val, 'f', -1, 64))

	return b.String()
}
