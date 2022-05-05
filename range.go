package rx

import (
	"context"

	"golang.org/x/exp/constraints"
)

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range[T constraints.Integer](low, high T) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		for v := low; v < high; v++ {
			if err := ctx.Err(); err != nil {
				sink.Error(err)
				return
			}

			sink.Next(v)
		}

		sink.Complete()
	}
}

// Iota creates an Observable that emits an infinite sequence of integers
// starting from init.
func Iota[T constraints.Integer](init T) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		for v := init; ; v++ {
			if err := ctx.Err(); err != nil {
				sink.Error(err)
				return
			}

			sink.Next(v)
		}
	}
}
