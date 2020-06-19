package rx

import (
	"context"
)

// Mutex creates an Observer that passes all emissions to sink in a mutually
// exclusive way.
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

// MutexContext creates an Observer that passes all emissions to sink in a
// mutually exclusive way while ctx is still active.
func MutexContext(ctx context.Context, sink Observer) Observer {
	cx := make(chan Observer, 1)
	cx <- sink
	return func(t Notification) {
		if sink, ok := <-cx; ok {
			if ctx.Err() != nil {
				close(cx)
				return
			}
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
