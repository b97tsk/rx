package rx

import (
	"context"
)

type fromObservablesObservable struct {
	Observables []Observable
}

func (obs fromObservablesObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	for _, obs := range obs.Observables {
		if ctx.Err() != nil {
			return Done()
		}
		sink.Next(obs)
	}
	sink.Complete()
	return Done()
}

// FromObservables creates an Observable that emits some Observables
// you specify as arguments, one after the other, and then completes.
func FromObservables(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return fromObservablesObservable{observables}.Subscribe
}
