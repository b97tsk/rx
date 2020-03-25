package rx

import (
	"context"
)

type takeWhileObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs takeWhileObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	var (
		sourceIndex = -1
		observer    Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if obs.Predicate(t.Value, sourceIndex) {
				sink(t)
				break
			}

			observer = NopObserver
			sink.Complete()

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, ctx.Cancel
}

// TakeWhile creates an Observable that emits values emitted by the source
// Observable so long as each value satisfies the given predicate, and then
// completes as soon as this predicate is not satisfied.
//
// TakeWhile takes values from the source only while they pass the condition
// given. When the first value does not satisfy, it completes.
func (Operators) TakeWhile(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return takeWhileObservable{source, predicate}.Subscribe
	}
}
