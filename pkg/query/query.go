package query

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/tmw/promgrep/pkg/tokenizer"
)

type Op string

const (
	OpEq       Op = "="
	OpNotEq       = "!="
	OpMatch       = "=~"
	OpNotMatch    = "!~"
)

type Matcher struct {
	Op  Op
	Val string
}

type Query struct {
	MetricName Matcher
	Labels     map[string]Matcher
}

func Compile(s string) (Query, error) {
	t := tokenizer.NewTokenizer(strings.NewReader(s), tokenizeText)
	return parse(t.Tokens())
}

func parse(t iter.Seq[Token]) (Query, error) {
	var (
		q = Query{
			Labels: make(map[string]Matcher),
		}

		tokens = slices.Collect(t)
	)

	for idx := 0; idx < len(tokens); {
		t := tokens[idx]
		if t.Typ == TokenTypeMetricName {
			q.MetricName = Matcher{Op: OpEq, Val: t.Str}
			idx++
			continue
		}

		if t.Typ == TokenTypeLabelName && t.Str == "__name__" {
			idx++
			op, numTokens := matchOperator(tokens[idx:])
			if numTokens == 0 {
				return q, fmt.Errorf("expected operator after __name__")
			}
			idx += numTokens

			valToken := tokens[idx]
			if valToken.Typ != TokenTypeLabelValue {
				return q, fmt.Errorf("expected label value after operator %s", op)
			}

			idx++
			q.MetricName = Matcher{Op: op, Val: valToken.Str}
			continue
		}

		if t.Typ == TokenTypeLabelName {
			idx++
			op, numTokens := matchOperator(tokens[idx:])
			if numTokens == 0 {
				return q, fmt.Errorf("expected operator after %s", t.Str)
			}
			idx += numTokens

			valToken := tokens[idx]
			if valToken.Typ != TokenTypeLabelValue {
				return q, fmt.Errorf("expected label value after operator %s", op)
			}

			q.Labels[t.Str] = Matcher{Op: op, Val: valToken.Str}
			idx++
			continue
		}

		idx++
	}

	return q, nil
}

// returns the operation it matched and the number of tokens it used
func matchOperator(tokens []Token) (Op, int) {
	if len(tokens) < 2 {
		panic("unable to parse operator; not enough tokens in buffer")
	}

	switch {
	case tokens[0].Typ == TokenTypeEq && tokens[1].Typ == TokenTypeTilde:
		return OpMatch, 2

	case tokens[0].Typ == TokenTypeExclamation && tokens[1].Typ == TokenTypeTilde:
		return OpNotMatch, 2

	case tokens[0].Typ == TokenTypeExclamation && tokens[1].Typ == TokenTypeEq:
		return OpNotEq, 2

	case tokens[0].Typ == TokenTypeEq: // why isn't it triggering on this?..
		return OpEq, 1
	}

	return OpEq, 0
}
