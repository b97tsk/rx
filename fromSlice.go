package rx

import (
	"context"
)

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice(slice []interface{}) Observable {
	return func(ctx context.Context, sink Observer) {
		for _, val := range slice {
			if ctx.Err() != nil {
				return
			}
			sink.Next(val)
		}
		sink.Complete()
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(values ...interface{}) Observable {
	return FromSlice(values)
}
