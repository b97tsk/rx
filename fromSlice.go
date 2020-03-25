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
	return func(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx := NewContext(parent)
		for _, val := range slice {
			if ctx.Err() != nil {
				return ctx, ctx.Cancel
			}
			sink.Next(val)
		}
		sink.Complete()
		ctx.Unsubscribe(Complete)
		return ctx, ctx.Cancel
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(values ...interface{}) Observable {
	return FromSlice(values)
}
