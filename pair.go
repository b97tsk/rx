package rx

import "golang.org/x/exp/constraints"

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

// FromMap creates an [Observable] that emits Pairs from a map, one after
// the other, and then completes. The order of those Pairs is not specified.
func FromMap[M ~map[K]V, K comparable, V any](m M) Observable[Pair[K, V]] {
	return func(c Context, o Observer[Pair[K, V]]) {
		done := c.Done()

		for k, v := range m {
			select {
			default:
			case <-done:
				o.Stop(c.Cause())
				return
			}

			Try1(o, Next(NewPair(k, v)), func() { o.Stop(ErrOops) })
		}

		o.Complete()
	}
}

// KeyOf maps each [Pair] emitted by the source [Observable] to a value stored
// in the Key field.
func KeyOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], K] {
	return LeftOf[Pair[K, V]]()
}

// LeftOf is an alias for [KeyOf].
func LeftOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], K] {
	return NewOperator(
		func(source Observable[Pair[K, V]]) Observable[K] {
			return func(c Context, o Observer[K]) {
				source.Subscribe(c, func(n Notification[Pair[K, V]]) {
					switch n.Kind {
					case KindNext:
						o.Next(n.Value.Left())
					case KindComplete:
						o.Complete()
					case KindError:
						o.Error(n.Error)
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}

// ValueOf maps each [Pair] emitted by the source [Observable] to a value
// stored in the Value field.
func ValueOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], V] {
	return RightOf[Pair[K, V]]()
}

// RightOf is an alias for [ValueOf].
func RightOf[_ Pair[K, V], K, V any]() Operator[Pair[K, V], V] {
	return NewOperator(
		func(source Observable[Pair[K, V]]) Observable[V] {
			return func(c Context, o Observer[V]) {
				source.Subscribe(c, func(n Notification[Pair[K, V]]) {
					switch n.Kind {
					case KindNext:
						o.Next(n.Value.Right())
					case KindComplete:
						o.Complete()
					case KindError:
						o.Error(n.Error)
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}

// Enumerate maps each value emitted by the source [Observable] to a [Pair]
// where the Key field stores the index of each value starting from init
// and the Value field stores each value.
func Enumerate[V any, K constraints.Integer](init K) Operator[V, Pair[K, V]] {
	return NewOperator(
		func(source Observable[V]) Observable[Pair[K, V]] {
			return func(c Context, o Observer[Pair[K, V]]) {
				i := init
				source.Subscribe(c, func(n Notification[V]) {
					switch n.Kind {
					case KindNext:
						o.Next(NewPair(i, n.Value))
						i++
					case KindComplete:
						o.Complete()
					case KindError:
						o.Error(n.Error)
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}
