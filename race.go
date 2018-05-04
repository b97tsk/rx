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
		mutable := MutableObserver{}
		mutable.Observer = ObserverFunc(func(t Notification) {
			if try.Lock() {
				for i, cancel := range subscriptions {
					if i != index {
						cancel()
					}
				}
				try.CancelAndUnlock()
				mutable.Observer = withFinalizer(ob, cancel)
				t.Observe(mutable.Observer)
			}
		})
		_, cancel := obsv.Subscribe(ctx, &mutable)
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
