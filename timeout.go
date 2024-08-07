package rx

import (
	"sync/atomic"
	"time"
)

// Timeout mirrors the source Observable, or emits a notification of ErrTimeout
// if the source does not emit a value in given time span.
func Timeout[T any](d time.Duration) TimeoutOperator[T] {
	return TimeoutOperator[T]{
		ts: timeoutConfig[T]{
			first: d,
			each:  d,
		},
	}
}

type timeoutConfig[T any] struct {
	first time.Duration
	each  time.Duration
	with  Observable[T]
}

// TimeoutOperator is an [Operator] type for [Timeout].
type TimeoutOperator[T any] struct {
	ts timeoutConfig[T]
}

// WithFirst sets First option to a given value.
func (op TimeoutOperator[T]) WithFirst(d time.Duration) TimeoutOperator[T] {
	op.ts.first = d
	return op
}

// WithObservable sets With option to a given value.
func (op TimeoutOperator[T]) WithObservable(ob Observable[T]) TimeoutOperator[T] {
	op.ts.with = ob
	return op
}

// Apply implements the Operator interface.
func (op TimeoutOperator[T]) Apply(source Observable[T]) Observable[T] {
	return timeoutObservable[T]{source, op.ts}.Subscribe
}

type timeoutObservable[T any] struct {
	source Observable[T]
	timeoutConfig[T]
}

func (ob timeoutObservable[T]) Subscribe(parent Context, o Observer[T]) {
	c, cancel := parent.WithCancel()

	var x struct {
		context atomic.Value
	}

	x.context.Store(c.Context)

	tm := time.AfterFunc(ob.first, c.PreAsyncCall(func() {
		if x.context.Swap(sentinel) != sentinel {
			cancel()

			if ob.with != nil {
				ob.with.Subscribe(parent, o)
				return
			}

			o.Error(ErrTimeout)
		}
	}))

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if tm.Stop() {
				Try1(o, n, func() {
					if c.WaitGroup != nil {
						c.WaitGroup.Done()
					}
				})
				tm.Reset(ob.each)
			}

		case KindError, KindComplete:
			if tm.Stop() {
				if c.WaitGroup != nil {
					c.WaitGroup.Done()
				}
			}

			if x.context.Swap(sentinel) != sentinel {
				cancel()
				o.Emit(n)
			}
		}
	})
}
