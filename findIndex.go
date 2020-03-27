package rx

import (
	"context"
)

type findIndexObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs findIndexObservable) Subscribe(ctx context.Context, sink Observer) {
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
				sink.Next(sourceIndex)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// FindIndex creates an Observable that emits only the index of the first value
// emitted by the source Observable that meets some condition.
func (Operators) FindIndex(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		obs := findIndexObservable{source, predicate}
		return Create(obs.Subscribe)
	}
}
