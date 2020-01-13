package rx

import (
	"context"
)

type raceObservable struct {
	Observables []Observable
}

func (obs raceObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Subscriptions []context.CancelFunc
	}
	cx := make(chan *X, 1)
	cx <- &X{
		Subscriptions: make([]context.CancelFunc, 0, len(obs.Observables)),
	}

	for index, obs := range obs.Observables {
		index := index

		x, ok := <-cx
		if !ok {
			break
		}

		ctx, cancel := context.WithCancel(ctx)
		x.Subscriptions = append(x.Subscriptions, cancel)

		cx <- x

		var observer Observer

		observer = func(t Notification) {
			if x, ok := <-cx; ok {
				for i, cancel := range x.Subscriptions {
					if i != index {
						cancel()
					}
				}
				close(cx)
				observer = sink
				sink(t)
				return
			}
			observer = NopObserver
			cancel()
		}

		obs.Subscribe(ctx, observer.Notify)
	}

	return ctx, cancel
}

// Race creates an Observable that mirrors the first source Observable to emit
// an item from the combination of this Observable and supplied Observables.
func Race(observables ...Observable) Observable {
	switch len(observables) {
	case 0:
		return Empty()
	case 1:
		return observables[0]
	default:
		return raceObservable{observables}.Subscribe
	}
}
