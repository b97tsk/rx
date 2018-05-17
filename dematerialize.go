package rx

import (
	"context"
)

type dematerializeOperator struct{}

func (op dematerializeOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			if t, ok := t.Value.(Notification); ok {
				switch {
				case t.HasValue:
					t.Observe(ob)
				default:
					mutableObserver = NopObserver
					t.Observe(ob)
					cancel()
				}
			} else {
				mutableObserver = NopObserver
				ob.Error(ErrNotNotification)
				cancel()
			}
		default:
			t.Observe(ob)
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent.
//
// Unwraps Notification objects as actual Next, Error and Complete emissions.
// The opposite of Materialize.
func (o Observable) Dematerialize() Observable {
	op := dematerializeOperator{}
	return o.Lift(op.Call)
}
