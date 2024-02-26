package rx

import (
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
// Internally, a runtime finalizer is set to call Error(ErrFinalized), which
// may run any time after Subject's Observer part gets garbage-collected.
func Multicast[T any]() Subject[T] {
	m := new(multicast[T])

	runtime.SetFinalizer(&m, multicastFinalizer[T])

	return Subject[T]{
		Observable: m.subscribe, // Only Observer field holds &m, this doesn't.
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

func (m *multicast[T]) subscribe(c Context, sink Observer[T]) {
	m.mu.Lock()

	last := m.last
	if last.Kind == 0 {
		var cancel CancelFunc

		c, cancel = c.WithCancel()
		sink = sink.OnLastNotification(cancel).Serialized()

		observer := sink
		m.obs.Add(&observer)

		c.AfterFunc(func() {
			m.mu.Lock()
			m.obs.Delete(&observer)
			m.mu.Unlock()
			observer.Error(c.Err())
		})
	}

	m.mu.Unlock()

	if last.Kind != 0 {
		sink(last)
	}
}
