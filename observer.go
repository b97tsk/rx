package rx

import "runtime"

// An Observer is a consumer of notifications delivered by an [Observable].
type Observer[T any] func(n Notification[T])

// Noop gives you an [Observer] that does nothing.
func Noop[T any](Notification[T]) {}

// NewObserver creates an [Observer] from f.
func NewObserver[T any](f func(n Notification[T])) Observer[T] { return f }

// Next passes a [Next] notification to o.
func (o Observer[T]) Next(v T) {
	o.Emit(Next(v))
}

// Complete passes a [Complete] notification to o.
func (o Observer[T]) Complete() {
	o.Emit(Complete[T]())
}

// Error passes an [Error] notification to o.
func (o Observer[T]) Error(err error) {
	o.Emit(Error[T](err))
}

// Stop passes a [Stop] notification to o.
func (o Observer[T]) Stop(err error) {
	o.Emit(Stop[T](err))
}

// Unsubscribe passes a [Stop] notification of ErrUnsubscribed to o.
func (o Observer[T]) Unsubscribe() {
	o.Emit(Stop[T](ErrUnsubscribed))
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

// DoOnTermination creates an [Observer] that passes incoming emissions to o,
// and when a notification of [Complete], [Error] or [Stop] passes in, calls
// f just before passing it to o.
//
// Note that [Stop] notifications may be emitted from random goroutines.
// If that happens, one would have to deal with race conditions.
// For more information, please refer to the package documentation.
func (o Observer[T]) DoOnTermination(f func()) Observer[T] {
	return func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			o.Emit(n)
		case KindComplete, KindError, KindStop:
			Try0(f, func() { o.Stop(ErrOops) })
			o.Emit(n)
		}
	}
}

// WithRuntimeFinalizer creates an Observer with a runtime finalizer set to
// run o.Stop(ErrFinalized) in a goroutine.
// o must be safe for concurrent use.
func WithRuntimeFinalizer[T any](o Observer[T]) Observer[T] {
	runtime.SetFinalizer(&o, func(o *Observer[T]) { go o.Stop(ErrFinalized) })
	return func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
		case KindComplete, KindError, KindStop:
			runtime.SetFinalizer(&o, nil)
		}
		o.Emit(n)
	}
}
