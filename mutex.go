package rx

import (
	"context"
)

// Mutex creates an Observer that passes all emissions to the specified
// Observer in a mutually exclusive way.
func Mutex(sink Observer) Observer {
	var try cancellableLocker
	return func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				sink(t)
				try.Unlock()
			default:
				try.CancelAndUnlock()
				sink(t)
			}
		}
	}
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func (Operators) Mutex() OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, Mutex(sink))
			},
		)
	}
}
