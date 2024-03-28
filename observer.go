package rx

import "runtime"

// An Observer is a consumer of notifications delivered by an [Observable].
type Observer[T any] func(n Notification[T])

// Noop gives you an Observer that does nothing.
func Noop[T any](Notification[T]) {}

// NewObserver creates an Observer from f.
func NewObserver[T any](f func(n Notification[T])) Observer[T] { return f }

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
	if n.Kind == KindNext {
		sink(n)
	}
}

// OnTermination creates an Observer that passes incoming emissions to sink,
// and when a notification of error or completion passes in, calls f just
// before passing it to sink.
func (sink Observer[T]) OnTermination(f func()) Observer[T] {
	return func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			sink(n)
		case KindError, KindComplete:
			Try0(f, func() { sink.Error(ErrOops) })
			sink(n)
		}
	}
}

// WithRuntimeFinalizer creates an Observer with a runtime finalizer set to
// run sink.Error(ErrFinalized) in a goroutine.
// sink must be safe for concurrent use.
func WithRuntimeFinalizer[T any](sink Observer[T]) Observer[T] {
	runtime.SetFinalizer(&sink, func(sink *Observer[T]) {
		go sink.Error(ErrFinalized)
	})

	return func(n Notification[T]) {
		switch n.Kind {
		case KindError, KindComplete:
			runtime.SetFinalizer(&sink, nil)
		}

		sink(n)
	}
}
