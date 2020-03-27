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
	return Create(
		func(ctx context.Context, sink Observer) {
			for _, obs := range observables {
				if ctx.Err() != nil {
					return
				}
				sink.Next(obs)
			}
			sink.Complete()
		},
	)
}
