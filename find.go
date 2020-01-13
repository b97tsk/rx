package rx

import (
	"context"
)

type findObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs findObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		outerIndex = -1
		observer   Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if obs.Predicate(t.Value, outerIndex) {
				observer = NopObserver
				sink(t)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (Operators) Find(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return findObservable{source, predicate}.Subscribe
	}
}
