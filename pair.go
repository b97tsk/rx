package rx

import (
	"context"

	"golang.org/x/exp/constraints"
)

// A Pair is a struct of two elements.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// Left returns p.Key.
func (p Pair[K, V]) Left() K { return p.Key }

// Right returns p.Value.
func (p Pair[K, V]) Right() V { return p.Value }

// MakePair creates a Pair of two elements.
func MakePair[K, V any](k K, v V) Pair[K, V] { return Pair[K, V]{k, v} }

// FromMap creates an Observable that emits Pairs from a map, one after
// the other, and then completes. The order of those Pairs is not specified.
func FromMap[M ~map[K]V, K comparable, V any](m M) Observable[Pair[K, V]] {
	return func(ctx context.Context, sink Observer[Pair[K, V]]) {
		for k, v := range m {
			if err := ctx.Err(); err != nil {
				sink.Error(err)
				return
			}

			sink.Next(MakePair(k, v))
		}

		sink.Complete()
	}
}
