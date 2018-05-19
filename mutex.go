package rx

import (
	"context"
)

type mutexOperator struct{}

func (op mutexOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var try cancellableLocker
	return source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				defer try.Unlock()
				sink(t)
			default:
				try.CancelAndUnlock()
				sink(t)
			}
		}
	})
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func (o Observable) Mutex() Observable {
	op := mutexOperator{}
	return o.Lift(op.Call)
}
