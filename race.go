package rx

import (
	"context"
)

type raceOperator struct {
	Observables []Observable
}

func (op raceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	length := len(op.Observables)
	subscriptions := make([]context.CancelFunc, 0, length)

	var try cancellableLocker

	for index, obsv := range op.Observables {
		index := index

		var observer Observer
		observer = func(t Notification) {
			if try.Lock() {
				for i, cancel := range subscriptions {
					if i != index {
						cancel()
					}
				}
				try.CancelAndUnlock()
				observer = Finally(sink, cancel)
				observer.Notify(t)
			}
		}
		_, cancel := obsv.Subscribe(ctx, observer.Notify)

		if try.Lock() {
			subscriptions = append(subscriptions, cancel)
			try.Unlock()
		} else {
			break
		}
	}

	return ctx, cancel
}

// Race creates an Observable that mirrors the first source Observable to emit
// an item from the combination of this Observable and supplied Observables.
func (Operators) Race(observables ...Observable) OperatorFunc {
	return func(source Observable) Observable {
		observables = append([]Observable{source}, observables...)
		op := raceOperator{observables}
		return Observable{}.Lift(op.Call)
	}
}
