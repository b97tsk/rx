package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type singleObservable struct {
	Source rx.Observable
}

func (obs singleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		value    interface{}
		hasValue bool
		observer rx.Observer
	)

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			if hasValue {
				observer = rx.Noop
				sink.Error(rx.ErrNotSingle)
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
				sink.Error(rx.ErrEmpty)
			}
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// notify of an ErrNotSingle or ErrEmpty respectively.
func Single() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := singleObservable{source}
		return rx.Create(obs.Subscribe)
	}
}
