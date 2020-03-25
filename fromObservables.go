package rx

import (
	"context"
)

// FromObservables creates an Observable that emits some Observables
// you specify as arguments, one after the other, and then completes.
func FromObservables(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return func(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx := NewContext(parent)
		for _, obs := range observables {
			if ctx.Err() != nil {
				return ctx, ctx.Cancel
			}
			sink.Next(obs)
		}
		sink.Complete()
		ctx.Unsubscribe(Complete)
		return ctx, ctx.Cancel
	}
}
