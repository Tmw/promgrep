package query

import (
	"fmt"
	"strings"

	"github.com/tmw/promgrep/pkg/tokenizer"
)

type TokenType string

const (
	TokenTypeMetricName  = "metric"
	TokenTypeLabelName   = "labelname"
	TokenTypeLabelValue  = "labelvalue"
	TokenTypeEq          = "eq"
	TokenTypeExclamation = "exclam"
	TokenTypeTilde       = "tilde"
)

type Token struct {
	Typ TokenType
	Str string
}

func Do() {
	str := `metric_name{cpu="12", color!="yellow", cluster=~"some-regex-.*", cluster_region!~"eu-*"}`
	// str := `{__name__=~"thing_.*"}`
	t := tokenizer.NewTokenizer(strings.NewReader(str), tokenizeText)
	for tok := range t.Tokens() {
		fmt.Println("tok", tok)
	}
}

func tokenizeText(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsOneOf(' '))

	if t.PeekMatch(tokenizer.IsEqual('{')) {
		return nil, tokenizeLabels
	}

	return nil, tokenizeMetricName
}

func tokenizeLabels(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsOneOf('{', ','))
	return nil, tokenizeLabelName
}

func tokenizeLabelName(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsEqual(' '))
	labelName := t.ReadUntil(tokenizer.IsOneOf('!', '='))
	tok := &Token{
		Typ: TokenTypeLabelName,
		Str: string(labelName),
	}

	return tok, tokenizeOperator
}

func tokenizeOperator(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	switch t.Peek() {
	case '!':
		return &Token{
			Typ: TokenTypeExclamation,
			Str: string(t.NextRune()),
		}, tokenizeOperator

	case '=':
		return &Token{
			Typ: TokenTypeEq,
			Str: string(t.NextRune()),
		}, tokenizeOperator

	case '~':
		return &Token{
			Typ: TokenTypeTilde,
			Str: string(t.NextRune()),
		}, tokenizeOperator

	case '"':
		return nil, tokenizeLabelValue
	}

	return nil, nil
}

func tokenizeLabelValue(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreUntil(tokenizer.IsEqual('"'))
	t.Ignore()

	var labelVal []rune

	// if not empty string, read value
	if !t.PeekMatch(tokenizer.IsEqual('"')) {
		labelVal = t.ReadUntil(tokenizer.IsEqual('"'))
		t.Ignore()
	}

	tok := &Token{
		Typ: TokenTypeLabelValue,
		Str: string(labelVal),
	}

	if t.PeekMatch(tokenizer.IsOneOf('}', '\n')) {
		return tok, nil
	}

	return tok, tokenizeLabels
}

func tokenizeMetricName(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	metricName := t.ReadUntil(tokenizer.IsOneOf('{', '\n'))
	tok := &Token{
		Typ: TokenTypeMetricName,
		Str: string(metricName),
	}

	if t.PeekMatch(tokenizer.IsEqual('{')) {
		return tok, tokenizeLabels
	}

	return tok, nil
}
