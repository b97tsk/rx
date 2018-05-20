package rx

import (
	"context"
)

type mutexOperator struct{}

func (op mutexOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, Mutex(sink))
}

// Mutex creates an Observer that passes all emissions to the specified
// Observer in a mutually exclusive way.
func Mutex(sink Observer) Observer {
	var try cancellableLocker
	return func(t Notification) {
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
	}
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func (o Observable) Mutex() Observable {
	op := mutexOperator{}
	return o.Lift(op.Call)
}
