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

// NewPair creates a Pair of two elements.
func NewPair[K, V any](k K, v V) Pair[K, V] { return Pair[K, V]{k, v} }

// FromMap creates an Observable that emits Pairs from a map, one after
// the other, and then completes. The order of those Pairs is not specified.
func FromMap[M ~map[K]V, K comparable, V any](m M) Observable[Pair[K, V]] {
	return func(ctx context.Context, sink Observer[Pair[K, V]]) {
		done := ctx.Done()

		for k, v := range m {
			select {
			default:
			case <-done:
				sink.Error(ctx.Err())
				return
			}

			sink.Next(NewPair(k, v))
		}

		sink.Complete()
	}
}

// KeyOf projects each Pair emitted by the source Observable to a value
// stored in the Key field of that Pair.
func KeyOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], K] {
	return LeftOf[Pair[K, V]]()
}

// LeftOf is an alias to KeyOf.
func LeftOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], K] {
	return NewOperator(
		func(source Observable[Pair[K, V]]) Observable[K] {
			return func(ctx context.Context, sink Observer[K]) {
				source.Subscribe(ctx, func(n Notification[Pair[K, V]]) {
					switch n.Kind {
					case KindNext:
						sink.Next(n.Value.Left())
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
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
	return RightOf[Pair[K, V]]()
}

// RightOf is an alias to ValueOf.
func RightOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], V] {
	return NewOperator(
		func(source Observable[Pair[K, V]]) Observable[V] {
			return func(ctx context.Context, sink Observer[V]) {
				source.Subscribe(ctx, func(n Notification[Pair[K, V]]) {
					switch n.Kind {
					case KindNext:
						sink.Next(n.Value.Right())
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
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
					switch n.Kind {
					case KindNext:
						sink.Next(NewPair(index, n.Value))
						index++
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						sink.Complete()
					}
				})
			}
		},
	)
}
