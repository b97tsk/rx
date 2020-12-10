package rx

import (
	"context"
)

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	return func(ctx context.Context, sink Observer) {
		for idx := low; idx < high; idx++ {
			if ctx.Err() != nil {
				return
			}

			sink.Next(idx)
		}

		sink.Complete()
	}
}
