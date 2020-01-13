package rx

import (
	"context"
)

// Mutex creates an Observer that passes all emissions to the specified
// Observer in a mutually exclusive way.
func Mutex(sink Observer) Observer {
	cx := make(chan Observer, 1)
	cx <- sink
	return func(t Notification) {
		if sink, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sink(t)
				cx <- sink
			default:
				close(cx)
				sink(t)
			}
		}
	}
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func (Operators) Mutex() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, Mutex(sink))
		}
	}
}
