package rx

import (
	"context"
)

type dematerializeOperator struct {
	source Operator
}

func (op dematerializeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if t, ok := t.Value.(Notification); ok {
				switch {
				case t.HasValue:
					ob.Next(t.Value)
				case t.HasError:
					mutable.Observer = NopObserver
					ob.Error(t.Value.(error))
					cancel()
				default:
					mutable.Observer = NopObserver
					ob.Complete()
					cancel()
				}
			} else {
				mutable.Observer = NopObserver
				ob.Error(ErrNotNotification)
				cancel()
			}
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			ob.Complete()
			cancel()
		}
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent.
//
// Unwraps Notification objects as actual Next, Error and Complete emissions.
// The opposite of Materialize.
func (o Observable) Dematerialize() Observable {
	op := dematerializeOperator{o.Op}
	return Observable{op}
}
