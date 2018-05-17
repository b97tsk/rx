package rx

import (
	"context"
)

type toObservablesOperator struct {
	source Operator
}

func (op toObservablesOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	observables := []Observable(nil)

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			if obsv, ok := t.Value.(Observable); ok {
				observables = append(observables, obsv)
			} else {
				mutableObserver = NopObserver
				ob.Error(ErrNotObservable)
				cancel()
			}
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			ob.Next(observables)
			ob.Complete()
			cancel()
		}
	}

	op.source.Call(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// ToObservables creates an Observable that collects all the Observables the
// source emits, then emits them as a slice of Observable when the source
// completes.
func (o Observable) ToObservables() Observable {
	op := toObservablesOperator{o.Op}
	return Observable{op}
}
