package rx

import (
	"time"

	"github.com/b97tsk/rx/internal/timerpool"
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
	c, cancel := parent.WithCancel()

	q := make(chan Notification[T])
	noop := make(chan struct{})

	c.Go(func() {
		tm := timerpool.Get(ob.First)

		for {
			select {
			case n := <-q:
				switch n.Kind {
				case KindNext:
					Try1(o, n, func() {
						cancel()
						close(noop)
						o.Error(ErrOops)
					})
				case KindError, KindComplete:
					timerpool.Put(tm)
					cancel()
					close(noop)
					o.Emit(n)
					return
				}

				tm.Reset(ob.Each)

			case <-tm.C:
				timerpool.PutExpired(tm)

				cancel()
				close(noop)

				if ob.With != nil {
					ob.With.Subscribe(parent, o)
					return
				}

				o.Error(ErrTimeout)

				return
			}
		}
	})

	ob.Source.Subscribe(c, channelObserver(q, noop))
}
