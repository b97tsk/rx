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

// KeyOf projects each Pair emitted by the source Observable to a value
// stored in the Key field of that Pair.
func KeyOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], K] {
	return NewOperator(
		func(source Observable[Pair[K, V]]) Observable[K] {
			return func(ctx context.Context, sink Observer[K]) {
				source.Subscribe(ctx, func(n Notification[Pair[K, V]]) {
					switch {
					case n.HasValue:
						sink.Next(n.Value.Left())
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Complete()
					}
				})
			}
		},
	)
}

// ValueOf projects each Pair emitted by the source Observable to a value
// stored in the Value field of that Pair.
func ValueOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], V] {
	return NewOperator(
		func(source Observable[Pair[K, V]]) Observable[V] {
			return func(ctx context.Context, sink Observer[V]) {
				source.Subscribe(ctx, func(n Notification[Pair[K, V]]) {
					switch {
					case n.HasValue:
						sink.Next(n.Value.Right())
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Complete()
					}
				})
			}
		},
	)
}

// WithIndex projects each value emitted by the source Observable to a Pair
// containing two elements: the Key field stores the index of each value
// starting from init; the Value field stores the value.
func WithIndex[V any, K constraints.Integer](init K) Operator[V, Pair[K, V]] {
	return NewOperator(
		func(source Observable[V]) Observable[Pair[K, V]] {
			return func(ctx context.Context, sink Observer[Pair[K, V]]) {
				index := init

				source.Subscribe(ctx, func(n Notification[V]) {
					switch {
					case n.HasValue:
						sink.Next(MakePair(index, n.Value))
						index++
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Complete()
					}
				})
			}
		},
	)
}
