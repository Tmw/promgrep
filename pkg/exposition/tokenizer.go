package exposition

import (
	"io"

	"github.com/tmw/promgrep/pkg/tokenizer"
)

func NewTokenizer(input io.RuneScanner) *tokenizer.Tokenizer[Token] {
	return tokenizer.NewTokenizer(input, tokenizeText)
}

func tokenizeText(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsEqual('\n'))

	if t.PeekMatch(tokenizer.IsEqual('#')) {
		return nil, tokenizeComment
	}

	if !t.PeekMatch(tokenizer.IsOneOf('#', '\n')) {
		return nil, tokenizeMetric
	}

	if t.Done() {
		return nil, nil
	}

	return nil, nil
}

func tokenizeMetric(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	name := t.ReadUntil(tokenizer.IsOneOf('{', ' '))
	tok := &Token{
		typ: TokenTypeMetric,
		str: string(name),
	}

	if t.PeekMatch(tokenizer.IsEqual('{')) {
		return tok, tokenizeLabelName
	}

	if t.PeekMatch(tokenizer.IsEqual(' ')) {
		return tok, tokenizeNumber
	}

	panic("should be unreachable")
}

func tokenizeLabelName(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsOneOf('{', ',', ' '))
	labelName := t.ReadUntil(tokenizer.IsEqual('='))

	return &Token{
		typ: TokenTypeLabelName,
		str: string(labelName),
	}, tokenizeLabelValue
}

func tokenizeLabelValue(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreUntil(tokenizer.IsEqual('"'))
	t.Ignore()

	var labelVal []rune

	// if not empty string, read value
	if !t.PeekMatch(tokenizer.IsEqual('"')) {
		labelVal = t.ReadUntil(tokenizer.IsEqual('"'))
	}

	tok := &Token{
		typ: TokenTypeLabelValue,
		str: string(labelVal),
	}

	t.IgnoreWhile(tokenizer.IsEqual('"'))

	if t.PeekMatch(tokenizer.IsEqual(',')) {
		return tok, tokenizeLabelName
	}

	if t.PeekMatch(tokenizer.IsEqual('}')) {
		return tok, tokenizeNumber
	}

	panic("should be unreachable")
}

func tokenizeNumber(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsOneOf('}', ' '))
	n := t.ReadUntil(tokenizer.IsEqual('\n'))

	tok := &Token{
		typ: TokenTypeNumber,
		str: string(n),
	}

	return tok, tokenizeText
}

func tokenizeComment(t *tokenizer.Tokenizer[Token]) (*Token, tokenizer.StateFn[Token]) {
	t.IgnoreWhile(tokenizer.IsOneOf(' ', '#'))

	typStr := string(t.ReadUntil(tokenizer.IsEqual(' ')))
	t.IgnoreWhile(tokenizer.IsEqual(' '))
	value := t.ReadUntil(tokenizer.IsEqual('\n'))

	if typStr == HELP {
		tok := &Token{
			typ: TokenTypeHelp,
			str: string(value),
		}

		return tok, tokenizeText
	}

	if typStr == TYPE {
		tok := &Token{
			typ: TokenTypeType,
			str: string(value),
		}

		return tok, tokenizeText
	}

	return nil, nil
}
