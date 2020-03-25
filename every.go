package rx

import (
	"context"
)

type everyObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs everyObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
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

			if !obs.Predicate(t.Value, sourceIndex) {
				observer = NopObserver
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

	return ctx, ctx.Cancel
}

// Every creates an Observable that emits whether or not every item of the source
// satisfies the condition specified.
//
// Every emits true or false, then completes.
func (Operators) Every(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return everyObservable{source, predicate}.Subscribe
	}
}
