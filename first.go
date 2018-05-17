package rx

import (
	"context"
)

type firstOperator struct{}

func (op firstOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			mutableObserver = NopObserver
			t.Observe(ob)
			ob.Complete()
			cancel()
		case t.HasError:
			t.Observe(ob)
			cancel()
		default:
			ob.Error(ErrEmpty)
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func (o Observable) First() Observable {
	op := firstOperator{}
	return o.Lift(op.Call)
}
