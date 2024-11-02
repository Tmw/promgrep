package maputil

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// Sorted is an iterator that, given a map, will iterate
// over its keys in a sorted fashion.
func Sorted[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	keys := slices.Sorted(maps.Keys(m))

	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				break
			}
		}
	}
}
