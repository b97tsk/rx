package rx

import (
	"context"
)

type raceOperator struct {
	observables []Observable
}

func (op raceOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	length := len(op.observables)
	subscriptions := make([]context.CancelFunc, 0, length)

	try := cancellableLocker{}

	for index, obsv := range op.observables {
		index := index

		var mutableObserver Observer

		mutableObserver = func(t Notification) {
			if try.Lock() {
				for i, cancel := range subscriptions {
					if i != index {
						cancel()
					}
				}
				try.CancelAndUnlock()
				mutableObserver = withFinalizer(ob, cancel)
				t.Observe(mutableObserver)
			}
		}

		_, cancel := obsv.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

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
func (o Observable) Race(observables ...Observable) Observable {
	observables = append([]Observable{o}, observables...)
	op := raceOperator{observables}
	return Observable{op}
}
