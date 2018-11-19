package rx

import (
	"context"
)

type raceOperator struct {
	Observables []Observable
}

func (op raceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	length := len(op.Observables)
	subscriptions := make([]context.CancelFunc, 0, length)

	var try cancellableLocker

	for index, obs := range op.Observables {
		index := index

		if !try.Lock() {
			break
		}

		ctx, cancel := context.WithCancel(ctx)
		subscriptions = append(subscriptions, cancel)

		try.Unlock()

		var observer Observer

		observer = func(t Notification) {
			if try.Lock() {
				for i, cancel := range subscriptions {
					if i != index {
						cancel()
					}
				}
				try.CancelAndUnlock()
				observer = sink
				sink(t)
				return
			}
			cancel()
		}

		obs.Subscribe(ctx, observer.Notify)
	}

	return ctx, cancel
}

// Race creates an Observable that mirrors the first source Observable to emit
// an item from the combination of this Observable and supplied Observables.
func Race(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	op := raceOperator{observables}
	return Observable{}.Lift(op.Call)
}
