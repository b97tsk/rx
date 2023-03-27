package rx

import (
	"context"
	"runtime"
	"sync"

	"github.com/b97tsk/rx/internal/ctxwatch"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// Multicast returns a Subject whose Observable part takes care of
// all Observers that subscribed to it, all of which will receive
// emissions from Subject's Observer part.
//
// Subjects are subject to memory leaks.
// After finishing using a Subject, you should call either its Error method
// or its Complete method to avoid that.
// If you can guarantee that every subscription to a Subject is canceled
// sooner or later, then you are fine.
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
			m.emit(n)

			if !n.HasValue {
				runtime.SetFinalizer(&m, nil)
			}
		},
	}
}

func multicastFinalizer[T any](m **multicast[T]) {
	go (*m).emit(Error[T](ErrFinalized))
}

type multicast[T any] struct {
	mu  sync.Mutex
	obs multiObserver[T]
	err error
}

func (m *multicast[T]) emit(n Notification[T]) {
	m.mu.Lock()

	switch {
	case m.err != nil:
		m.mu.Unlock()

	case n.HasValue:
		obs := m.obs.Clone()
		defer obs.Release()

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Emit(n)
		}

	default:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.err = errComplete
		if n.HasError {
			m.err = errOrErrNil(n.Error)
		}

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Emit(n)
		}
	}
}

func (m *multicast[T]) subscribe(ctx context.Context, sink Observer[T]) {
	m.mu.Lock()

	err := m.err
	if err == nil {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.OnLastNotification(cancel).Serialized()

		observer := sink
		m.obs.Add(&observer)

		wg := waitgroup.Get(ctx)
		if wg != nil {
			wg.Add(1)
		}

		ctxwatch.Add(ctx, func() {
			if wg != nil {
				defer wg.Done()
			}

			m.mu.Lock()
			m.obs.Delete(&observer)
			m.mu.Unlock()

			observer.Error(ctx.Err())
		})
	}

	m.mu.Unlock()

	if err != nil {
		if err == errComplete {
			sink.Complete()
			return
		}

		sink.Error(cleanErrNil(err))
	}
}
