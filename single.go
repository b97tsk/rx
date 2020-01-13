package rx

import (
	"context"
)

type singleObservable struct {
	Source Observable
}

func (obs singleObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		value    interface{}
		hasValue bool
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if hasValue {
				observer = NopObserver
				sink.Error(ErrNotSingle)
			} else {
				value = t.Value
				hasValue = true
			}
		case t.HasError:
			sink(t)
		default:
			if hasValue {
				sink.Next(value)
				sink.Complete()
			} else {
				sink.Error(ErrEmpty)
			}
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// notify of an ErrNotSingle or ErrEmpty respectively.
func (Operators) Single() OperatorFunc {
	return func(source Observable) Observable {
		return singleObservable{source}.Subscribe
	}
}
