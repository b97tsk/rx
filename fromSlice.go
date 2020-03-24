package rx

import (
	"context"
)

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice(slice []interface{}) Observable {
	if len(slice) == 0 {
		return Empty()
	}
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		for _, val := range slice {
			if ctx.Err() != nil {
				return Done(ctx)
			}
			sink.Next(val)
		}
		sink.Complete()
		return Done(ctx)
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(values ...interface{}) Observable {
	return FromSlice(values)
}
