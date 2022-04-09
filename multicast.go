package rx

import (
	"context"
	"runtime"
	"sync"

	"github.com/b97tsk/rx/internal/ctxwatch"
)

// Multicast returns a Subject whose Observable part takes care of all
// Observers that subscribes to it, which will receive emissions from
// Subject's Observer part.
func Multicast() Subject {
	m := &multicast{}

	runtime.SetFinalizer(&m, multicastFinalizer)

	return Subject{
		Observable: func(ctx context.Context, sink Observer) {
			m.subscribe(ctx, sink)
		},
		Observer: func(t Notification) {
			m.sink(t)
		},
	}
}

func multicastFinalizer(m **multicast) {
	go Observer((*m).sink).Error(errFinalized)
}

type multicast struct {
	mu  sync.Mutex
	err error
	lst observerList
}

func (m *multicast) sink(t Notification) {
	m.mu.Lock()

	switch {
	case m.err != nil:
		m.mu.Unlock()

	case t.HasValue:
		lst := m.lst.Clone()
		defer lst.Release()

		m.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList

		m.lst.Swap(&lst)

		m.err = errCompleted

		if t.HasError {
			m.err = t.Error

			if m.err == nil {
				m.err = errNil
			}
		}

		m.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (m *multicast) subscribe(ctx context.Context, sink Observer) {
	m.mu.Lock()

	err := m.err
	if err == nil {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.WithCancel(cancel).MutexContext(ctx)

		observer := sink
		m.lst.Append(&observer)
		ctxwatch.Add(ctx, func() {
			m.mu.Lock()
			m.lst.Remove(&observer)
			m.mu.Unlock()
		})
	}

	m.mu.Unlock()

	if err != nil {
		if err == errCompleted {
			sink.Complete()
			return
		}

		if err == errNil {
			err = nil
		}

		sink.Error(err)
	}
}
