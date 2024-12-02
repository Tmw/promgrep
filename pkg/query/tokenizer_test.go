package query

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmw/promgrep/pkg/tokenizer"
)

func TestTokenizer(t *testing.T) {
	type testcase struct {
		input    string
		expected []Token
	}

	cases := map[string]testcase{
		"metric with single label/value pair": {
			input: "metric label=value",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric"},
				{Typ: TokenTypeLabelName, Str: "label"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value"},
			},
		},

		"metric with two label/value pairs": {
			input: "metric label_a=value_a label_b=value_b",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},
		},

		"no metric with label/value pair": {
			input: "label_a=value_a",
			expected: []Token{
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
			},
		},

		"no metric with multiple label/value pairs": {
			input: "label_a=value_a label_b=value_b",
			expected: []Token{
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},
		},

		"just metric name label": {
			input: "__name__=some_metric",
			expected: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "some_metric"},
			},
		},

		"with metric name label and label/value pair": {
			input: "__name__=some_metric label_a=value_a",
			expected: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "some_metric"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
			},
		},

		"with metric and negating label/value pair": {
			input: "metric_name label!=value",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value"},
			},
		},

		"with metric and negating label/value pair and standard label/value pair": {
			input: "metric_name label_a!=value_a label_b=value_b",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},
		},

		"with metric and regex label/value pair": {
			input: "metric_name label_a=~value_.+",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_.+"},
			},
		},

		"with metric and negating regex label/value pair": {
			input: "metric_name label_a!~value_.+",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_.+"},
			},
		},

		"with metric name label and regex label/value pair": {
			input: "__name__=~value_.+",
			expected: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_.+"},
			},
		},

		"with metric name label and negating regex label/value pair": {
			input: "__name__!~value_.+",
			expected: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_.+"},
			},
		},

		"just metric name": {
			input: "metric_name",
			expected: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
			},
		},

		"empty string": {
			input:    "",
			expected: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			stream := tokenizer.NewTokenizer(strings.NewReader(tc.input), tokenizeText)
			assert.Equal(t, tc.expected, slices.Collect(stream.Tokens()))
		})
	}
}
