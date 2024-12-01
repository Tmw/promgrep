package query

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_metric_name(t *testing.T) {
	type testcase struct {
		input    []Token
		expected Query
	}

	cases := map[string]testcase{
		"just metric name": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels:     make(map[string]Matcher),
			},
		},

		"metric name label (eq)": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "metric_name"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels:     make(map[string]Matcher),
			},
		},

		"metric name label (neq)": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "metric_name"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpNotEq, Val: "metric_name"},
				Labels:     make(map[string]Matcher),
			},
		},

		"metric name label (match)": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "metric_.+"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpMatch, Val: "metric_.+"},
				Labels:     make(map[string]Matcher),
			},
		},

		"metric name label (negated match)": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "metric_.+"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpNotMatch, Val: "metric_.+"},
				Labels:     make(map[string]Matcher),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := parse(slices.Values(tc.input))
			assert.NoError(t, err, "error parsing query")
			assert.Equal(t, tc.expected, actual, "expected query")
		})
	}
}

func TestQuery_labels(t *testing.T) {
	type testcase struct {
		input    []Token
		expected Query
	}

	cases := map[string]testcase{
		"single label (eq)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpEq, Val: "value_a"},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := parse(slices.Values(tc.input))
			assert.NoError(t, err, "error parsing query")
			assert.Equal(t, tc.expected, actual, "expected query")
		})
	}
}
