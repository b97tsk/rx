package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/timerpool"
)

// Timeout mirrors the source Observable, throws ErrTimeout if the source
// does not emit a value in given time span.
func Timeout[T any](d time.Duration) TimeoutOperator[T] {
	return TimeoutOperator[T]{
		opts: timeoutConfig[T]{
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

// TimeoutOperator is an Operator type for Timeout.
type TimeoutOperator[T any] struct {
	opts timeoutConfig[T]
}

// WithFirst sets First option to a given value.
func (op TimeoutOperator[T]) WithFirst(d time.Duration) TimeoutOperator[T] {
	op.opts.First = d
	return op
}

// WithObservable sets With option to a given value.
func (op TimeoutOperator[T]) WithObservable(obs Observable[T]) TimeoutOperator[T] {
	op.opts.With = obs
	return op
}

// Apply implements the Operator interface.
func (op TimeoutOperator[T]) Apply(source Observable[T]) Observable[T] {
	return timeoutObservable[T]{source, op.opts}.Subscribe
}

// AsOperator converts op to an Operator.
//
// Once type inference has improved in Go, this method will be removed.
func (op TimeoutOperator[T]) AsOperator() Operator[T, T] { return op }

type timeoutObservable[T any] struct {
	Source Observable[T]
	timeoutConfig[T]
}

func (obs timeoutObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	childCtx, cancel := context.WithCancel(ctx)

	c := make(chan Notification[T])
	noop := make(chan struct{})

	Go(ctx, func() {
		tm := timerpool.Get(obs.First)

		for {
			select {
			case n := <-c:
				sink(n)

				if !n.HasValue {
					timerpool.Put(tm)

					cancel()
					close(noop)

					return
				}

				tm.Reset(obs.Each)

			case <-tm.C:
				timerpool.PutExpired(tm)

				cancel()
				close(noop)

				if obs.With != nil {
					obs.With.Subscribe(ctx, sink)
					return
				}

				sink.Error(ErrTimeout)

				return
			}
		}
	})

	obs.Source.Subscribe(childCtx, chanObserver(c, noop))
}
