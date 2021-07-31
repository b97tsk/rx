package rx

import (
	"context"
)

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	return func(ctx context.Context, sink Observer) {
		for val := low; val < high; val++ {
			if ctx.Err() != nil {
				return
			}

			sink.Next(val)
		}

		sink.Complete()
	}
}

// Iota creates an Observable that emits an infinite sequence of integers
// starting from init.
func Iota(init int) Observable {
	return func(ctx context.Context, sink Observer) {
		for val := init; ; val++ {
			if ctx.Err() != nil {
				return
			}

			sink.Next(val)
		}
	}
}
