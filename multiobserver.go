package rx

import (
	"github.com/b97tsk/rx/internal/atomic"
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

		rc.Sub(1)
	}
}

func (m *multiObserver[T]) Add(observer *Observer[T]) {
	observers := m.Observers
	oldcap := cap(observers)

	observers = append(observers, observer)
	m.Observers = observers

	if cap(observers) != oldcap {
		if rc := m.rc; rc != nil && !rc.Equal(0) {
			m.rc = nil
		}
	}
}

func (m *multiObserver[T]) Delete(observer *Observer[T]) {
	observers := m.Observers

	for i, sink := range observers {
		if sink == observer {
			n := len(observers)

			if rc := m.rc; rc != nil && !rc.Equal(0) {
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

			break
		}
	}
}
