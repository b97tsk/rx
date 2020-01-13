package rx

import (
	"context"
)

type onErrorResumeNextObservable struct {
	Observables []Observable
}

func (obs onErrorResumeNextObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		observer       Observer
		avoidRecursive avoidRecursiveCalls
	)

	remainder := obs.Observables
	subscribe := func() {
		source := remainder[0]
		remainder = remainder[1:]
		source.Subscribe(ctx, observer)
	}

	observer = func(t Notification) {
		if t.HasValue {
			sink(t)
			return
		}
		if len(remainder) == 0 {
			sink.Complete()
			return
		}
		avoidRecursive.Do(subscribe)
	}

	avoidRecursive.Do(subscribe)

	return ctx, cancel
}

// OnErrorResumeNext creates an Observable that concatenates the provided
// Observables. It's like Concat, but errors are ignored.
func OnErrorResumeNext(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return onErrorResumeNextObservable{observables}.Subscribe
}
