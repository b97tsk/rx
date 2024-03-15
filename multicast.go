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
	Mu    sync.Mutex
	Mobs  multiObserver[T]
	LastN Notification[struct{}]
}

func (m *multicast[T]) emit(n Notification[T]) {
	m.Mu.Lock()

	if m.LastN.Kind != 0 {
		m.Mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
		mobs := m.Mobs.Clone()
		defer mobs.Release()

		m.Mu.Unlock()

		mobs.Emit(n)

	case KindError, KindComplete:
		var mobs multiObserver[T]

		m.Mobs, mobs = mobs, m.Mobs

		switch n.Kind {
		case KindError:
			m.LastN = Error[struct{}](n.Error)
		case KindComplete:
			m.LastN = Complete[struct{}]()
		}

		m.Mu.Unlock()

		mobs.Emit(n)

	default: // Unknown kind.
		m.Mu.Unlock()
	}
}

func (m *multicast[T]) subscribe(c Context, sink Observer[T]) {
	m.Mu.Lock()

	lastn := m.LastN
	if lastn.Kind == 0 {
		var cancel CancelFunc

		c, cancel = c.WithCancel()
		sink = sink.OnLastNotification(cancel).Serialized()

		observer := sink
		m.Mobs.Add(&observer)

		c.AfterFunc(func() {
			m.Mu.Lock()
			m.Mobs.Delete(&observer)
			m.Mu.Unlock()
			observer.Error(c.Err())
		})
	}

	m.Mu.Unlock()

	switch lastn.Kind {
	case KindError:
		sink.Error(lastn.Error)
	case KindComplete:
		sink.Complete()
	}
}
