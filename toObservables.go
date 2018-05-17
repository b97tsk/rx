package rx

import (
	"context"
)

type toObservablesOperator struct{}

func (op toObservablesOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		observables     []Observable
		mutableObserver Observer
	)

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
			t.Observe(ob)
			cancel()
		default:
			ob.Next(observables)
			ob.Complete()
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// ToObservables creates an Observable that collects all the Observables the
// source emits, then emits them as a slice of Observable when the source
// completes.
func (o Observable) ToObservables() Observable {
	op := toObservablesOperator{}
	return o.Lift(op.Call)
}
