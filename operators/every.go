package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type everyObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs everyObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		sourceIndex = -1
		observer    rx.Observer
	)

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if !obs.Predicate(t.Value, sourceIndex) {
				observer = rx.NopObserver
				sink.Next(false)
				sink.Complete()
			}

		case t.HasError:
			sink(t)

		default:
			sink.Next(true)
			sink.Complete()
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// Every creates an Observable that emits whether or not every item of the source
// satisfies the condition specified.
//
// Every emits true or false, then completes.
func Every(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := everyObservable{source, predicate}
		return rx.Create(obs.Subscribe)
	}
}
