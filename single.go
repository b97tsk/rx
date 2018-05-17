package rx

import (
	"context"
)

type singleOperator struct{}

func (op singleOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		value           interface{}
		hasValue        bool
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			if hasValue {
				mutableObserver = NopObserver
				ob.Error(ErrNotSingle)
				cancel()
			} else {
				value = t.Value
				hasValue = true
			}
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			if hasValue {
				ob.Next(value)
				ob.Complete()
			} else {
				ob.Error(ErrEmpty)
			}
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// notify of an ErrNotSingle or ErrEmpty respectively.
func (o Observable) Single() Observable {
	op := singleOperator{}
	return o.Lift(op.Call)
}
