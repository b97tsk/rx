package rx

import (
	"sync/atomic"
)

type multiObserver[T any] struct {
	observers []*Observer[T]
	refcount  *atomic.Uint32
}

func (m *multiObserver[T]) Clone() multiObserver[T] {
	rc := m.refcount

	if rc == nil {
		rc = new(atomic.Uint32)
		m.refcount = rc
	}

	rc.Add(1)

	return multiObserver[T]{m.observers, rc}
}

func (m *multiObserver[T]) Release() {
	if rc := m.refcount; rc != nil {
		m.refcount = nil
		rc.Add(^uint32(0))
	}
}

func (m *multiObserver[T]) Add(o *Observer[T]) {
	observers := m.observers
	oldcap := cap(observers)

	observers = append(observers, o)
	m.observers = observers

	if cap(observers) != oldcap {
		if rc := m.refcount; rc != nil && rc.Load() != 0 {
			m.refcount = nil
		}
	}
}

func (m *multiObserver[T]) Delete(o *Observer[T]) {
	observers := m.observers

	for i := range observers {
		if o == observers[i] {
			n := len(observers)

			if rc := m.refcount; rc != nil && rc.Load() != 0 {
				new := make([]*Observer[T], n-1, n)
				copy(new, observers[:i])
				copy(new[i:], observers[i+1:])
				m.observers = new
				m.refcount = nil
			} else {
				copy(observers[i:], observers[i+1:])
				observers[n-1] = nil
				m.observers = observers[:n-1]
			}

			return
		}
	}
}

func (m *multiObserver[T]) Emit(n Notification[T]) {
	emitNotificationToObservers(m.observers, n)
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
