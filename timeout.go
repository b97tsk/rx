package rx

import (
	"time"

	"github.com/b97tsk/rx/internal/timerpool"
)

// Timeout mirrors the source Observable, or emits an error notification
// of ErrTimeout if the source does not emit a value in given time span.
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

// TimeoutOperator is an [Operator] type for [Timeout].
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

type timeoutObservable[T any] struct {
	Source Observable[T]
	timeoutConfig[T]
}

func (obs timeoutObservable[T]) Subscribe(parent Context, sink Observer[T]) {
	c, cancel := parent.WithCancel()

	q := make(chan Notification[T])
	noop := make(chan struct{})

	c.Go(func() {
		tm := timerpool.Get(obs.First)

		for {
			select {
			case n := <-q:
				switch n.Kind {
				case KindError, KindComplete:
					cancel()
					close(noop)
				}

				sink(n)

				switch n.Kind {
				case KindError, KindComplete:
					timerpool.Put(tm)
					return
				}

				tm.Reset(obs.Each)

			case <-tm.C:
				timerpool.PutExpired(tm)

				cancel()
				close(noop)

				if obs.With != nil {
					obs.With.Subscribe(parent, sink)
					return
				}

				sink.Error(ErrTimeout)

				return
			}
		}
	})

	obs.Source.Subscribe(c, channelObserver(q, noop))
}
