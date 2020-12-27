package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/critical"
)

// Race creates an Observable that mirrors the first source Observable to emit
// an item from the combination of this Observable and supplied Observables.
func Race(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}

	return raceObservable(observables).Subscribe
}

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

	var race critical.Section

	for i, obs := range observables {
		var observer Observer

		index := i

		observer = func(t Notification) {
			if critical.Enter(&race) {
				for i := range subscriptions {
					if i != index {
						subscriptions[i].Cancel()
					}
				}

				critical.Close(&race)

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
