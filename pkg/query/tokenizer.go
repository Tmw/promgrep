package query

import (
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

func tokenizeText(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsOneOf(' '))

	ident := t.ReadUntil(tokenizer.IsOneOf(' ', '=', '!', '\n'))
	if t.PeekMatch(tokenizer.IsEqual(' ')) {
		return &Token{
			Typ: TokenTypeMetricName,
			Str: string(ident),
		}, tokenizeText
	}

	if t.PeekMatch(tokenizer.IsOneOf('=', '!')) {
		tok := &Token{
			Typ: TokenTypeLabelName,
			Str: string(ident),
		}

		switch t.Peek() {
		case '=':
			return tok, tokenizeEq
		case '!':
			return tok, tokenizeExclam
		}
	}

	return nil, nil
}

func tokenizeEq(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.Ignore()
	tok := &Token{
		Typ: TokenTypeEq,
		Str: "=",
	}

	if t.PeekMatch(tokenizer.IsEqual('~')) {
		return tok, tokenizeTilde
	}

	return tok, tokenizeLabelValue
}

func tokenizeTilde(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.Ignore()
	return &Token{
		Typ: TokenTypeTilde,
		Str: "~",
	}, tokenizeLabelValue
}

func tokenizeExclam(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.Ignore()

	tok := &Token{
		Typ: TokenTypeExclamation,
		Str: "!",
	}

	if t.PeekMatch(tokenizer.IsEqual('~')) {
		return tok, tokenizeTilde
	}

	if t.PeekMatch(tokenizer.IsEqual('=')) {
		return tok, tokenizeEq
	}

	return nil, nil
}

func tokenizeLabelValue(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	labelVal := t.ReadUntil(tokenizer.IsOneOf(' ', '\n'))
	return &Token{
		Typ: TokenTypeLabelValue,
		Str: string(labelVal),
	}, tokenizeText
}
