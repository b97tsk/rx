package rx

import "runtime"

// An Observer is a consumer of notifications delivered by an [Observable].
type Observer[T any] func(n Notification[T])

// Noop gives you an Observer that does nothing.
func Noop[T any](Notification[T]) {}

// NewObserver creates an Observer from f.
func NewObserver[T any](f func(n Notification[T])) Observer[T] { return f }

// Next passes a value to o.
func (o Observer[T]) Next(v T) {
	o.Emit(Next(v))
}

// Error passes an error to o.
func (o Observer[T]) Error(err error) {
	o.Emit(Error[T](err))
}

// Complete passes a completion to o.
func (o Observer[T]) Complete() {
	o.Emit(Complete[T]())
}

// Emit passes n to o.
func (o Observer[T]) Emit(n Notification[T]) {
	o(n)
}

// ElementsOnly passes n to o if n represents a value.
func (o Observer[T]) ElementsOnly(n Notification[T]) {
	if n.Kind == KindNext {
		o.Emit(n)
	}
}

// DoOnTermination creates an Observer that passes incoming emissions to o,
// and when a notification of error or completion passes in, calls f just
// before passing it to o.
func (o Observer[T]) DoOnTermination(f func()) Observer[T] {
	return func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			o.Emit(n)
		case KindError, KindComplete:
			Try0(f, func() { o.Error(ErrOops) })
			o.Emit(n)
		}
	}
}

// WithRuntimeFinalizer creates an Observer with a runtime finalizer set to
// run o.Error(ErrFinalized) in a goroutine.
// o must be safe for concurrent use.
func WithRuntimeFinalizer[T any](o Observer[T]) Observer[T] {
	runtime.SetFinalizer(&o, func(o *Observer[T]) {
		go o.Error(ErrFinalized)
	})

	return func(n Notification[T]) {
		switch n.Kind {
		case KindError, KindComplete:
			runtime.SetFinalizer(&o, nil)
		}

		o.Emit(n)
	}
}
