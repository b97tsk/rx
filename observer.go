package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/critical"
)

// An Observer is a consumer of notifications delivered by an Observable.
type Observer func(Notification)

// Next passes a value to sink.
func (sink Observer) Next(val interface{}) {
	sink(Notification{Value: val, HasValue: true})
}

// Error passes an error to sink.
func (sink Observer) Error(err error) {
	sink(Notification{Error: err, HasError: true})
}

// Complete passes a completion to sink.
func (sink Observer) Complete() {
	sink(Notification{})
}

// Sink passes t to *sink.
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

// ElementsOnly creates an Observer that only passes values to sink.
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
	var s critical.Section
	return func(t Notification) {
		if critical.Enter(&s) {
			switch {
			case t.HasValue:
				defer critical.Leave(&s)
				sink(t)
			default:
				critical.Close(&s)
				sink(t)
			}
		}
	}
}

// MutexContext creates an Observer that passes incoming emissions to sink in
// a mutually exclusive way while ctx is active.
func (sink Observer) MutexContext(ctx context.Context) Observer {
	var s critical.Section
	return func(t Notification) {
		if critical.Enter(&s) {
			switch {
			case ctx.Err() != nil:
				critical.Close(&s)
			case t.HasValue:
				defer critical.Leave(&s)
				sink(t)
			default:
				critical.Close(&s)
				sink(t)
			}
		}
	}
}

// WithCancel creates an Observer that passes incoming emissions to sink, and
// when an error or a completion passes in, calls a specified function just
// before passing it to sink.
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
