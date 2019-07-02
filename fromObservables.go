package rx

import (
	"context"
)

type fromObservablesOperator struct {
	Observables []Observable
}

func (op fromObservablesOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	for _, obs := range op.Observables {
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
	switch {
	case len(observables) > 1:
		op := fromObservablesOperator{observables}
		return Observable{}.Lift(op.Call)
	case len(observables) == 1:
		return just(observables[0])
	default:
		return Empty()
	}
}
