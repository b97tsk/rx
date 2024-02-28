package rx

import (
	"sync/atomic"
)

type multiObserver[T any] struct {
	Observers []*Observer[T]

	rc *atomic.Uint32
}

func (m *multiObserver[T]) Clone() multiObserver[T] {
	rc := m.rc

	if rc == nil {
		rc = new(atomic.Uint32)
		m.rc = rc
	}

	rc.Add(1)

	return multiObserver[T]{m.Observers, rc}
}

func (m *multiObserver[T]) Release() {
	if rc := m.rc; rc != nil {
		m.rc = nil

		rc.Add(^uint32(0))
	}
}

func (m *multiObserver[T]) Add(observer *Observer[T]) {
	observers := m.Observers
	oldcap := cap(observers)

	observers = append(observers, observer)
	m.Observers = observers

	if cap(observers) != oldcap {
		if rc := m.rc; rc != nil && rc.Load() != 0 {
			m.rc = nil
		}
	}
}

func (m *multiObserver[T]) Delete(observer *Observer[T]) {
	observers := m.Observers

	for i, sink := range observers {
		if sink == observer {
			n := len(observers)

			if rc := m.rc; rc != nil && rc.Load() != 0 {
				new := make([]*Observer[T], n-1, n)
				copy(new, observers[:i])
				copy(new[i:], observers[i+1:])
				m.Observers = new
				m.rc = nil
			} else {
				copy(observers[i:], observers[i+1:])
				observers[n-1] = nil
				m.Observers = observers[:n-1]
			}

			return
		}
	}
}

func (m *multiObserver[T]) Emit(n Notification[T]) {
	emitNotificationToObservers(m.Observers, n)
}

func emitNotificationToObservers[T any](observers []*Observer[T], n Notification[T]) {
	var i int

	defer func() {
		if i < len(observers)-1 {
			emitNotificationToObservers(observers[i+1:], n)
		}
	}()

	var sink *Observer[T]

	for i, sink = range observers {
		sink.Emit(n)
	}
}
