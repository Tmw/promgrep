package query

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tmw/promgrep/pkg/tokenizer"
)

func TestTokenizer(t *testing.T) {
	str := "metric label=value"
	stream := tokenizer.NewTokenizer(strings.NewReader(str), tokenizeText)
	for tok := range stream.Tokens() {
		fmt.Println(" - tok: ", tok)
	}

	t.Fail()
}
