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
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		for _, obs := range observables {
			if ctx.Err() != nil {
				return Done()
			}
			sink.Next(obs)
		}
		sink.Complete()
		return Done()
	}
}
