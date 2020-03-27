package rx

import (
	"context"
)

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	if low >= high {
		return Empty()
	}
	return Create(
		func(ctx context.Context, sink Observer) {
			for index := low; index < high; index++ {
				if ctx.Err() != nil {
					return
				}
				sink.Next(index)
			}
			sink.Complete()
		},
	)
}
