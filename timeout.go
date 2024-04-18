package rx

import (
	"context"
	"errors"
	"time"
)

// Timeout mirrors the source Observable, or emits a notification of ErrTimeout
// if the source does not emit a value in given time span.
func Timeout[T any](d time.Duration) TimeoutOperator[T] {
	return TimeoutOperator[T]{
		ts: timeoutConfig[T]{
			First: d,
			Each:  d,
		},
	}
}

type timeoutConfig[T any] struct {
	First time.Duration
	Each  time.Duration
	With  Observable[T]
}

// TimeoutOperator is an [Operator] type for [Timeout].
type TimeoutOperator[T any] struct {
	ts timeoutConfig[T]
}

// WithFirst sets First option to a given value.
func (op TimeoutOperator[T]) WithFirst(d time.Duration) TimeoutOperator[T] {
	op.ts.First = d
	return op
}

// WithObservable sets With option to a given value.
func (op TimeoutOperator[T]) WithObservable(ob Observable[T]) TimeoutOperator[T] {
	op.ts.With = ob
	return op
}

// Apply implements the Operator interface.
func (op TimeoutOperator[T]) Apply(source Observable[T]) Observable[T] {
	return timeoutObservable[T]{source, op.ts}.Subscribe
}

type timeoutObservable[T any] struct {
	Source Observable[T]
	timeoutConfig[T]
}

func (ob timeoutObservable[T]) Subscribe(parent Context, o Observer[T]) {
	c, cancel := parent.WithCancelCause()
	o = o.DoOnTermination(func() { cancel(nil) })

	timeout := errors.Join(context.Canceled)
	tm := time.AfterFunc(ob.First, func() { cancel(timeout) })

	ob.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if tm.Stop() {
				o.Emit(n)
				tm.Reset(ob.Each)
			}

		case KindError:
			tm.Stop()

			if errors.Is(n.Error, timeout) {
				if ob.With != nil {
					ob.With.Subscribe(parent, o)
					return
				}

				o.Error(ErrTimeout)

				return
			}

			o.Emit(n)

		case KindComplete:
			tm.Stop()
			o.Emit(n)
		}
	})
}
