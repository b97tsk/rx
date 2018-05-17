package rx

import (
	"context"
)

type singleOperator struct {
	source Operator
}

func (op singleOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	value := interface{}(nil)
	hasValue := false

	var mutableObserver Observer

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

	op.source.Call(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// notify of an ErrNotSingle or ErrEmpty respectively.
func (o Observable) Single() Observable {
	op := singleOperator{o.Op}
	return Observable{op}
}
