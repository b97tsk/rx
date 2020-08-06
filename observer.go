package rx

import (
	"context"
)

// An Observer is a consumer of notifications delivered by an Observable.
type Observer func(Notification)

// Next passes a NEXT emission to sink.
func (sink Observer) Next(val interface{}) {
	sink(Notification{Value: val, HasValue: true})
}

// Error passes an ERROR emission to sink.
func (sink Observer) Error(err error) {
	sink(Notification{Error: err, HasError: true})
}

// Complete passes a COMPLETE emission to sink.
func (sink Observer) Complete() {
	sink(Notification{})
}

// Sink passes a specified emission to *sink.
//
// Sink also yields an Observer that is equivalent to:
//
//	func(t Notification) { (*sink)(t) }
//
// Useful when you want to set *sink to another Observer at some point.
//
// Another use case: when you name your Observer observer, it looks bad if you
// call it like this: observer(t). Better use observer.Sink(t) instead.
//
func (sink *Observer) Sink(t Notification) {
	(*sink)(t)
}

// ElementsOnly creates an Observer that only passes NEXT emissions to sink.
func (sink Observer) ElementsOnly() Observer {
	return func(t Notification) {
		if t.HasValue {
			sink(t)
		}
	}
}

// Mutex creates an Observer that passes incoming emissions to sink in a
// mutually exclusive way.
func (sink Observer) Mutex() Observer {
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

// MutexContext creates an Observer that passes incoming emissions to sink in
// a mutually exclusive way while ctx is still active.
func (sink Observer) MutexContext(ctx context.Context) Observer {
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

// WithCancel creates an Observer that passes incoming emissions to sink and,
// when an ERROR or COMPLETE emission passes in, calls a specified function
// just before passing it to sink.
func (sink Observer) WithCancel(cancel context.CancelFunc) Observer {
	return func(t Notification) {
		if !t.HasValue {
			cancel()
		}
		sink(t)
	}
}

// Noop gives you an Observer that does nothing.
func Noop(Notification) {}
