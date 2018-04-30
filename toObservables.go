package rx

import (
	"context"
)

type toObservablesOperator struct {
	source Operator
}

func (op toObservablesOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	all := []Observable(nil)

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if obsv, ok := t.Value.(Observable); ok {
				all = append(all, obsv)
			} else {
				mutable.Observer = NopObserver
				ob.Error(ErrNotObservable)
				cancel()
			}
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			ob.Next(all)
			ob.Complete()
			cancel()
		}
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// ToObservables creates an Observable that collects all the Observables the
// source emits, then emits them as a slice of Observable when the source
// completes.
func (o Observable) ToObservables() Observable {
	op := toObservablesOperator{o.Op}
	return Observable{op}
}
