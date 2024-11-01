package tokenizer

import (
	"slices"
	"strconv"
)

func IsEqual(a rune) matchFn {
	return func(i rune) bool {
		return a == i
	}
}

func IsOneOf(a ...rune) matchFn {
	return func(i rune) bool {
		return slices.Contains(a, i)
	}
}

func IsNumeric(i rune) bool {
	_, err := strconv.ParseFloat(string(i), 64)
	return err == nil
}
