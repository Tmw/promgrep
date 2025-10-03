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

		"single label (not eq)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpNotEq, Val: "value_a"},
				},
			},
		},

		"single label (match)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpMatch, Val: "value_a"},
				},
			},
		},

		"single label (not match)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpNotMatch, Val: "value_a"},
				},
			},
		},

		"two labels (eq)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpEq, Val: "value_a"},
					"label_b": {Op: OpEq, Val: "value_b"},
				},
			},
		},

		"two labels (eq and match)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpEq, Val: "value_a"},
					"label_b": {Op: OpMatch, Val: "value_b"},
				},
			},
		},

		"two labels (eq and not match)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpEq, Val: "value_a"},
					"label_b": {Op: OpNotMatch, Val: "value_b"},
				},
			},
		},

		"two labels (not eq and not match)": {
			input: []Token{
				{Typ: TokenTypeMetricName, Str: "metric_name"},
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeLabelValue, Str: "value_a"},
				{Typ: TokenTypeLabelName, Str: "label_b"},
				{Typ: TokenTypeExclamation, Str: "!"},
				{Typ: TokenTypeTilde, Str: "~"},
				{Typ: TokenTypeLabelValue, Str: "value_b"},
			},

			expected: Query{
				MetricName: Matcher{Op: OpEq, Val: "metric_name"},
				Labels: map[string]Matcher{
					"label_a": {Op: OpNotEq, Val: "value_a"},
					"label_b": {Op: OpNotMatch, Val: "value_b"},
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

func TestQuery_unhappy_path(t *testing.T) {
	type testcase struct {
		input    []Token
		expected string
	}

	cases := map[string]testcase{
		"starts with operator (eq)": {
			input: []Token{
				{Typ: TokenTypeEq, Str: "="},
			},

			expected: "unexpected token eq",
		},

		"starts with operator (match)": {
			input: []Token{
				{Typ: TokenTypeTilde, Str: "~"},
			},

			expected: "unexpected token tilde",
		},

		"missing operator after __name__": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "__name__"},
			},

			expected: "expected operator after __name__",
		},

		"missing operator after label": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "label_a"},
			},

			expected: "expected operator after label_a",
		},

		"missing label value after operator": {
			input: []Token{
				{Typ: TokenTypeLabelName, Str: "label_a"},
				{Typ: TokenTypeEq, Str: "="},
			},

			expected: "expected label value after operator =",
		},

		"wrong operator combination": {
			input: []Token{
				{Typ: TokenTypeEq, Str: "="},
				{Typ: TokenTypeEq, Str: "="},
			},

			expected: "unexpected token eq",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := parse(slices.Values(tc.input))
			assert.ErrorContains(t, err, tc.expected)
		})
	}
}
