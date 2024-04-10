package rx

import (
	"sync/atomic"
)

type multiObserver[T any] struct {
	Observers []*Observer[T]
	RefCount  *atomic.Uint32
}

func (m *multiObserver[T]) Clone() multiObserver[T] {
	rc := m.RefCount

	if rc == nil {
		rc = new(atomic.Uint32)
		m.RefCount = rc
	}

	rc.Add(1)

	return multiObserver[T]{m.Observers, rc}
}

func (m *multiObserver[T]) Release() {
	if rc := m.RefCount; rc != nil {
		m.RefCount = nil
		rc.Add(^uint32(0))
	}
}

func (m *multiObserver[T]) Add(o *Observer[T]) {
	observers := m.Observers
	oldcap := cap(observers)

	observers = append(observers, o)
	m.Observers = observers

	if cap(observers) != oldcap {
		if rc := m.RefCount; rc != nil && rc.Load() != 0 {
			m.RefCount = nil
		}
	}
}

func (m *multiObserver[T]) Delete(o *Observer[T]) {
	observers := m.Observers

	for i := range observers {
		if o == observers[i] {
			n := len(observers)

			if rc := m.RefCount; rc != nil && rc.Load() != 0 {
				new := make([]*Observer[T], n-1, n)
				copy(new, observers[:i])
				copy(new[i:], observers[i+1:])
				m.Observers = new
				m.RefCount = nil
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

	var o *Observer[T]

	for i, o = range observers {
		o.Emit(n)
	}
}
