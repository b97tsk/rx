package rx

import (
	"context"
	"runtime"
	"sync"

	"github.com/b97tsk/rx/internal/ctxwatch"
)

// Multicast returns a Subject whose Observable part takes care of all
// Observers that subscribed to it, all of which will receive emissions
// from Subject's Observer part.
//
// Subjects are subject to memory leaks. Once you have done using a Subject,
// you should at least call its Error or Complete method once to avoid that.
// If you can guarantee that every subscription to the Subject is canceled,
// then you are fine.
//
// Internally, a runtime finalizer is set to call Error(ErrFinalized), but
// it is not guaranteed to work.
func Multicast[T any]() Subject[T] {
	m := new(multicast[T])

	runtime.SetFinalizer(&m, multicastFinalizer[T])

	return Subject[T]{
		Observable: func(ctx context.Context, sink Observer[T]) {
			m.subscribe(ctx, sink)
		},
		Observer: func(n Notification[T]) {
			m.sink(n)
		},
	}
}

func multicastFinalizer[T any](m **multicast[T]) {
	go (*m).sink(Error[T](ErrFinalized))
}

type multicast[T any] struct {
	mu  sync.Mutex
	obs multiObserver[T]
	err error
}

func (m *multicast[T]) sink(n Notification[T]) {
	m.mu.Lock()

	switch {
	case m.err != nil:
		m.mu.Unlock()

	case n.HasValue:
		obs := m.obs.Clone()
		defer obs.Release()

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Sink(n)
		}

	default:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.err = errCompleted
		if n.HasError {
			m.err = errOrErrNil(n.Error)
		}

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Sink(n)
		}
	}
}

func (m *multicast[T]) subscribe(ctx context.Context, sink Observer[T]) {
	m.mu.Lock()

	err := m.err
	if err == nil {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.WithCancel(cancel).Mutex()

		observer := sink
		m.obs.Add(&observer)
		ctxwatch.Add(ctx, func(ctx context.Context) {
			m.mu.Lock()
			m.obs.Delete(&observer)
			m.mu.Unlock()
			observer.Error(ctx.Err())
		})
	}

	m.mu.Unlock()

	if err != nil {
		if err == errCompleted {
			sink.Complete()
			return
		}

		sink.Error(cleanErrNil(err))
	}
}
