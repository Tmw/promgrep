package tokenizer

import (
	"io"
	"iter"
)

type matchFn func(c rune) bool
type StateFn[Token any] func(t *Tokenizer[Token]) (*Token, StateFn[Token])

type Tokenizer[Token any] struct {
	input io.RuneScanner
	done  bool
	state StateFn[Token]
}

func NewTokenizer[Token any](input io.RuneScanner, initialState StateFn[Token]) *Tokenizer[Token] {
	return &Tokenizer[Token]{
		input: input,
		done:  false,
		state: initialState,
	}
}

func (t *Tokenizer[Token]) Done() bool {
	return t.done
}

func (t *Tokenizer[Token]) NextRune() rune {
	r, _, err := t.input.ReadRune()

	if err == io.EOF {
		t.done = true
	}

	return r
}

func (t *Tokenizer[Token]) Peek() rune {
	r := t.NextRune()
	t.input.UnreadRune()

	return r
}

func (t *Tokenizer[Token]) PeekMatch(fn matchFn) bool {
	return fn(t.Peek())
}

func (t *Tokenizer[Token]) ReadUntil(fn matchFn) []rune {
	out := []rune{}

	for !t.done {
		// not efficient, but something to figure out later...
		out = append(out, t.NextRune())
		if fn(t.Peek()) {
			break
		}
	}

	return out
}

func (t *Tokenizer[Token]) IgnoreWhile(fn matchFn) {
	for !t.done && fn(t.Peek()) {
		t.Ignore()
	}
}

func (t *Tokenizer[Token]) IgnoreUntil(fn matchFn) {
	for !t.done && !fn(t.Peek()) {
		t.Ignore()
	}
}

func (t *Tokenizer[Token]) Ignore() {
	_ = t.NextRune()
}

func (t *Tokenizer[Token]) Tokens() iter.Seq[Token] {
	var tok *Token

	return func(yield func(Token) bool) {
		for !t.done && t.state != nil {
			tok, t.state = t.state(t)
			if tok != nil {
				if !yield(*tok) {
					break
				}
			}
		}
	}
}
