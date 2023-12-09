package rx

import (
	"context"
	"runtime"
	"sync"
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

			switch n.Kind {
			case KindError, KindComplete:
				runtime.SetFinalizer(&m, nil)
			}
		},
	}
}

func multicastFinalizer[T any](m **multicast[T]) {
	go (*m).emit(Error[T](ErrFinalized))
}

type multicast[T any] struct {
	mu   sync.Mutex
	obs  multiObserver[T]
	last Notification[T]
}

func (m *multicast[T]) emit(n Notification[T]) {
	m.mu.Lock()

	if m.last.Kind != 0 {
		m.mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
		obs := m.obs.Clone()
		defer obs.Release()

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Emit(n)
		}

	case KindError, KindComplete:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.last = n

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Emit(n)
		}

	default: // Unknown kind.
		m.mu.Unlock()
	}
}

func (m *multicast[T]) subscribe(ctx context.Context, sink Observer[T]) {
	m.mu.Lock()

	last := m.last
	if last.Kind == 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.OnLastNotification(cancel).Serialized()

		observer := sink
		m.obs.Add(&observer)

		wg := WaitGroupFromContext(ctx)
		if wg != nil {
			wg.Add(1)
		}

		context.AfterFunc(ctx, func() {
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

	if last.Kind != 0 {
		sink(last)
	}
}
