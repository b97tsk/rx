package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type findIndexObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs findIndexObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if obs.Predicate(t.Value, sourceIndex) {
				observer = rx.Noop
				sink.Next(sourceIndex)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// FindIndex creates an Observable that emits only the index of the first value
// emitted by the source Observable that meets some condition.
func FindIndex(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := findIndexObservable{source, predicate}
		return rx.Create(obs.Subscribe)
	}
}
