package exposition

import (
	"io"
	"iter"
	"slices"
	"strconv"
)

type matchFn func(c rune) bool
type stateFn func(t *Tokenizer) (*Token, stateFn)

type Tokenizer struct {
	input io.RuneScanner
	done  bool
}

func NewTokenizer(input io.RuneScanner) *Tokenizer {
	return &Tokenizer{
		input: input,
		done:  false,
	}
}

func (t *Tokenizer) nextRune() rune {
	r, _, err := t.input.ReadRune()

	if err == io.EOF {
		t.done = true
	}

	return r
}

func (t *Tokenizer) peek() rune {
	r := t.nextRune()
	t.input.UnreadRune()

	return r
}

func (t *Tokenizer) peekMatch(fn matchFn) bool {
	return fn(t.peek())
}

func (t *Tokenizer) readUntil(fn matchFn) []rune {
	out := []rune{}

	for !t.done {
		// not efficient, but something to figure out later...
		out = append(out, t.nextRune())
		if fn(t.peek()) {
			break
		}
	}

	return out
}

func (t *Tokenizer) ignoreWhile(fn matchFn) {
	for !t.done && fn(t.peek()) {
		t.ignore()
	}
}

func (t *Tokenizer) ignoreUntil(fn matchFn) {
	for !t.done && !fn(t.peek()) {
		t.ignore()
	}
}

func (t *Tokenizer) ignore() {
	_ = t.nextRune()
}

func isEqual(a rune) matchFn {
	return func(i rune) bool {
		return a == i
	}
}

func isOneOf(a ...rune) matchFn {
	return func(i rune) bool {
		return slices.Contains(a, i)
	}
}

func isNumeric(i rune) bool {
	_, err := strconv.ParseFloat(string(i), 64)
	return err == nil
}

func tokenizeText(t *Tokenizer) (*Token, stateFn) {
	t.ignoreWhile(isEqual('\n'))

	if t.peekMatch(isEqual('#')) {
		return nil, tokenizeComment
	}

	if !t.peekMatch(isOneOf('#', '\n')) {
		return nil, tokenizeMetric
	}

	if t.done {
		return nil, nil
	}

	return nil, nil
}

func tokenizeMetric(t *Tokenizer) (*Token, stateFn) {
	name := t.readUntil(isOneOf('{', ' '))
	tok := &Token{
		typ: TokenTypeMetric,
		str: string(name),
	}

	if t.peekMatch(isEqual('{')) {
		return tok, tokenizeLabelName
	}

	if t.peekMatch(isEqual(' ')) {
		return tok, tokenizeNumber
	}

	panic("should be unreachable")
}

func tokenizeLabelName(t *Tokenizer) (*Token, stateFn) {
	t.ignoreWhile(isOneOf('{', ',', ' '))
	labelName := t.readUntil(isEqual('='))

	return &Token{
		typ: TokenTypeLabelName,
		str: string(labelName),
	}, tokenizeLabelValue
}

func tokenizeLabelValue(t *Tokenizer) (*Token, stateFn) {
	t.ignoreUntil(isEqual('"'))
	t.ignore()

	var labelVal []rune

	// if not empty string, read value
	if !t.peekMatch(isEqual('"')) {
		labelVal = t.readUntil(isEqual('"'))
	}

	tok := &Token{
		typ: TokenTypeLabelValue,
		str: string(labelVal),
	}

	t.ignoreWhile(isEqual('"'))

	if t.peekMatch(isEqual(',')) {
		return tok, tokenizeLabelName
	}

	if t.peekMatch(isEqual('}')) {
		return tok, tokenizeNumber
	}

	panic("should be unreachable")
}

func tokenizeNumber(t *Tokenizer) (*Token, stateFn) {
	t.ignoreWhile(isOneOf('}', ' '))
	n := t.readUntil(isEqual('\n'))

	tok := &Token{
		typ: TokenTypeNumber,
		str: string(n),
	}

	return tok, tokenizeText
}

func tokenizeComment(t *Tokenizer) (*Token, stateFn) {
	t.ignoreWhile(isOneOf(' ', '#'))

	typStr := string(t.readUntil(isEqual(' ')))
	t.ignoreWhile(isEqual(' '))
	value := t.readUntil(isEqual('\n'))

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

func (t *Tokenizer) Tokens() iter.Seq[Token] {
	var (
		tok   *Token
		state stateFn = tokenizeText
	)

	return func(yield func(Token) bool) {
		for !t.done && state != nil {
			tok, state = state(t)
			if tok != nil {
				if !yield(*tok) {
					break
				}
			}
		}
	}
}
