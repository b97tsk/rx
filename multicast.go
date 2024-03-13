package rx

import "sync"

// Multicast returns a Subject that mirrors every emission it receives to all
// its subscribers.
func Multicast[T any]() Subject[T] {
	m := new(multicast[T])
	return Subject[T]{
		Observable: NewObservable(m.subscribe),
		Observer:   NewObserver(m.emit).WithRuntimeFinalizer(),
	}
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

		obs.Emit(n)

	case KindError, KindComplete:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.last = n

		m.mu.Unlock()

		obs.Emit(n)

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
