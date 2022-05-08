package rx

import (
	"context"
)

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice[S ~[]T, T any](s S) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		for _, v := range s {
			if err := ctx.Err(); err != nil {
				sink.Error(err)
				return
			}

			sink.Next(v)
		}

		sink.Complete()
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just[T any](s ...T) Observable[T] {
	return FromSlice(s)
}
