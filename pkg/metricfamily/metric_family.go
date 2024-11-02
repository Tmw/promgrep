package metricfamily

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tmw/promgrep/pkg/maputil"
)

type MetricFamily struct {
	Name   string
	Val    float64
	Labels map[string]string
}

func (m *MetricFamily) String() string {
	var b strings.Builder
	b.WriteString(m.Name)

	numLabels := len(m.Labels)
	curLabel := 0

	if numLabels > 0 {
		b.WriteRune('{')
		for k, v := range maputil.Sorted(m.Labels) {
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
