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
			op, numTokens, err := matchOperator(tokens[idx:])
			if err != nil {
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
			op, numTokens, err := matchOperator(tokens[idx:])
			if err != nil {
				return q, fmt.Errorf("expected operator after %s", t.Str)
			}

			idx += numTokens
			if idx > len(tokens)-1 {
				return q, fmt.Errorf("expected label value after operator %s", op)
			}

			valToken := tokens[idx]
			if valToken.Typ != TokenTypeLabelValue {
				return q, fmt.Errorf("expected label value after operator %s", op)
			}

			q.Labels[t.Str] = Matcher{Op: op, Val: valToken.Str}
			idx++
			continue
		}

		return q, fmt.Errorf("unexpected token %s", t.Typ)
	}

	return q, nil
}

// returns the operation it matched and the number of tokens it used
func matchOperator(tokens []Token) (Op, int, error) {
	matches := func(t []Token, idx int, typ TokenType) bool {
		if idx > len(t)-1 {
			return false
		}

		tok := t[idx]
		return tok.Typ == typ
	}

	switch {
	case matches(tokens, 0, TokenTypeEq) && matches(tokens, 1, TokenTypeTilde):
		return OpMatch, 2, nil

	case matches(tokens, 0, TokenTypeExclamation) && matches(tokens, 1, TokenTypeTilde):
		return OpNotMatch, 2, nil

	case matches(tokens, 0, TokenTypeExclamation) && matches(tokens, 1, TokenTypeEq):
		return OpNotEq, 2, nil

	case matches(tokens, 0, TokenTypeEq):
		return OpEq, 1, nil
	}

	return OpEq, 0, fmt.Errorf("unable to match operator")
}
