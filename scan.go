package rx

import (
	"context"
)

// Scan applies an accumulator function over the source Observable,
// and emits each intermediate result, given an initial value.
func Scan[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	if accumulator == nil {
		panic("accumulator == nil")
	}

	return scan(init, accumulator)
}

func scan[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(ctx context.Context, sink Observer[R]) {
				res := init

				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						res = accumulator(res, n.Value)

						sink.Next(res)
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
