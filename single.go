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

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if hasValue {
				mutable.Observer = NopObserver
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
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// notify of an ErrNotSingle or ErrEmpty respectively.
func (o Observable) Single() Observable {
	op := singleOperator{o.Op}
	return Observable{op}
}
