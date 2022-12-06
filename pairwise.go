package rx

import (
	"context"
)

// PairWise groups pairs of consecutive emissions together and emits them
// as Pairs.
func PairWise[T any]() Operator[T, Pair[T, T]] {
	return NewOperator(pairWise[T])
}

func pairWise[T any](source Observable[T]) Observable[Pair[T, T]] {
	return func(ctx context.Context, sink Observer[Pair[T, T]]) {
		var p struct {
			Value    T
			HasValue bool
		}

		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case n.HasValue:
				if p.HasValue {
					sink.Next(NewPair(p.Value, n.Value))
				}

				p.Value = n.Value
				p.HasValue = true

			case n.HasError:
				sink.Error(n.Error)

			default:
				sink.Complete()
			}
		})
	}
}
