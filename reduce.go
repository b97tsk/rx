package rx

import (
	"context"
)

// Reduce applies an accumulator function over the source Observable,
// and emits the accumulated result when the source completes, given
// an initial value.
func Reduce[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	if accumulator == nil {
		panic("accumulator == nil")
	}

	return reduce(init, accumulator)
}

func reduce[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	return AsOperator(
		func(source Observable[T]) Observable[R] {
			return func(ctx context.Context, sink Observer[R]) {
				res := init

				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						res = accumulator(res, n.Value)
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Next(res)
						sink.Complete()
					}
				})
			}
		},
	)
}
