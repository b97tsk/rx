package rx

import (
	"github.com/b97tsk/rx/internal/critical"
)

// An Observer is a consumer of notifications delivered by an [Observable].
type Observer[T any] func(n Notification[T])

// Next passes a value to sink.
func (sink Observer[T]) Next(v T) {
	sink(Next(v))
}

// Error passes an error to sink.
func (sink Observer[T]) Error(e error) {
	sink(Error[T](e))
}

// Complete passes a completion to sink.
func (sink Observer[T]) Complete() {
	sink(Complete[T]())
}

// Emit passes n to sink.
func (sink Observer[T]) Emit(n Notification[T]) {
	sink(n)
}

// ElementsOnly passes n to sink if n represents a value.
func (sink Observer[T]) ElementsOnly(n Notification[T]) {
	if n.HasValue {
		sink(n)
	}
}

// OnLastNotification creates an Observer that passes incoming emissions to
// sink, and when a notification of error or completion passes in, calls f
// just before passing it to sink.
func (sink Observer[T]) OnLastNotification(f func()) Observer[T] {
	return func(n Notification[T]) {
		if !n.HasValue {
			f()
		}

		sink(n)
	}
}

// WithMutex creates an Observer that passes incoming emissions to sink
// in a mutually exclusive way.
func (sink Observer[T]) WithMutex() Observer[T] {
	var lock critical.Section

	return func(n Notification[T]) {
		if critical.Enter(&lock) {
			switch {
			case n.HasValue:
				sink(n)
				critical.Leave(&lock)
			default:
				critical.Close(&lock)
				sink(n)
			}
		}
	}
}

// Noop gives you an Observer that does nothing.
func Noop[T any](Notification[T]) {}

// NewObserver creates an Observer from f.
func NewObserver[T any](f func(n Notification[T])) Observer[T] { return f }
