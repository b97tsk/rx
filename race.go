package rx

import (
	"context"
)

type raceObservable []Observable

func (observables raceObservable) Subscribe(ctx context.Context, sink Observer) {
	type Subscription struct {
		Context context.Context
		Cancel  context.CancelFunc
	}

	subscriptions := make([]Subscription, len(observables))
	for i := range subscriptions {
		ctx, cancel := context.WithCancel(ctx)
		subscriptions[i] = Subscription{ctx, cancel}
	}

	race := make(chan struct{}, 1)
	race <- struct{}{}

	for i, obs := range observables {
		var (
			index    = i
			observer Observer
		)
		observer = func(t Notification) {
			if _, ok := <-race; ok {
				for i := range subscriptions {
					if i != index {
						subscriptions[i].Cancel()
					}
				}
				close(race)
				observer = sink
				sink(t)
				return
			}
			observer = Noop
			subscriptions[index].Cancel()
		}
		go obs.Subscribe(subscriptions[i].Context, observer.Sink)
	}
}

// Race creates an Observable that mirrors the first source Observable to emit
// an item from the combination of this Observable and supplied Observables.
func Race(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return raceObservable(observables).Subscribe
}
