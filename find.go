package rx

import (
	"context"
)

type findObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs findObservable) Subscribe(ctx context.Context, sink Observer) {
	var (
		sourceIndex = -1
		observer    Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if obs.Predicate(t.Value, sourceIndex) {
				observer = NopObserver
				sink(t)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (Operators) Find(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		obs := findObservable{source, predicate}
		return Create(obs.Subscribe)
	}
}
