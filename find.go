package rx

import (
	"context"
)

type findObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs findObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
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
				observer = NopObserver
				sink(t)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, ctx.Cancel
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (Operators) Find(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return findObservable{source, predicate}.Subscribe
	}
}
