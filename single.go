package rx

import (
	"context"
)

type singleOperator struct{}

func (op singleOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

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
				cancel()
			} else {
				value = t.Value
				hasValue = true
			}
		case t.HasError:
			sink(t)
			cancel()
		default:
			if hasValue {
				sink.Next(value)
				sink.Complete()
			} else {
				sink.Error(ErrEmpty)
			}
			cancel()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// notify of an ErrNotSingle or ErrEmpty respectively.
func (o Observable) Single() Observable {
	op := singleOperator{}
	return o.Lift(op.Call)
}
