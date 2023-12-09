package rx

import (
	"context"
)

// Pairwise groups pairs of consecutive emissions together and emits them
// as Pairs.
func Pairwise[T any]() Operator[T, Pair[T, T]] {
	return NewOperator(pairwise[T])
}

func pairwise[T any](source Observable[T]) Observable[Pair[T, T]] {
	return func(ctx context.Context, sink Observer[Pair[T, T]]) {
		var p struct {
			Value    T
			HasValue bool
		}

		source.Subscribe(ctx, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				if p.HasValue {
					sink.Next(NewPair(p.Value, n.Value))
				}

				p.Value = n.Value
				p.HasValue = true

			case KindError:
				sink.Error(n.Error)

			case KindComplete:
				sink.Complete()
			}
		})
	}
}
