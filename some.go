package rx

import (
	"context"
)

type someObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs someObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
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
				sink.Next(true)
				sink.Complete()
			}

		case t.HasError:
			sink(t)

		default:
			sink.Next(false)
			sink.Complete()
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, ctx.Cancel
}

// Some creates an Observable that emits whether or not any item of the source
// satisfies the condition specified.
//
// Some emits true or false, then completes.
func (Operators) Some(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return someObservable{source, predicate}.Subscribe
	}
}
