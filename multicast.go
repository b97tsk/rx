package rx

import (
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
)

// Multicast returns a Subject that forwards every value it receives to
// all its subscribers.
// Values emitted to a Multicast before the first subscriber are lost.
func Multicast[T any]() Subject[T] {
	return MulticastBuffer[T](0)
}

// MulticastBufferAll returns a Subject that keeps track of every value
// it receives.
// Each subscriber will then receive all tracked values as well as future
// values.
func MulticastBufferAll[T any]() Subject[T] {
	return MulticastBuffer[T](-1)
}

// MulticastBuffer returns a Subject that keeps track of a certain number of
// recent values it receive.
// Each subscriber will then receive all tracked values as well as future
// values.
//
// If n < 0, MulticastBuffer keeps track of every value it receives;
// if n == 0, MulticastBuffer doesn't keep track of any value it receives
// at all.
func MulticastBuffer[T any](n int) Subject[T] {
	m := &multicast[T]{Cap: n}
	return Subject[T]{
		Observable: NewObservable(m.Subscribe),
		Observer:   WithRuntimeFinalizer(m.Emit),
	}
}

type multicast[T any] struct {
	Mu    sync.Mutex
	Cap   int
	Mobs  multiObserver[T]
	LastN Notification[struct{}]
	Buf   *struct {
		Queue    queue.Queue[T]
		RefCount atomic.Uint32
	}
}

func pnew[T any](*T) *T { return new(T) }

func (m *multicast[T]) Emit(n Notification[T]) {
	m.Mu.Lock()

	if m.LastN.Kind != 0 {
		m.Mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
		mobs := m.Mobs.Clone()
		defer mobs.Release()

		if m.Cap != 0 {
			b := m.Buf

			switch {
			case b == nil:
				b = pnew(b)
				m.Buf = b
			case b.RefCount.Load() != 0:
				q := b.Queue.Clone()
				b = pnew(b)
				b.Queue = q
				m.Buf = b
			}

			q := &b.Queue

			if q.Len() == m.Cap {
				q.Pop()
			}

			q.Push(n.Value)
		}

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

func (m *multicast[T]) Subscribe(c Context, sink Observer[T]) {
	m.Mu.Lock()

	lastn := m.LastN
	if lastn.Kind == 0 {
		c, sink = Serialize(c, sink)

		observer := sink
		m.Mobs.Add(&observer)

		c.AfterFunc(func() {
			m.Mu.Lock()
			m.Mobs.Delete(&observer)
			m.Mu.Unlock()
			observer.Error(c.Err())
		})
	}

	b := m.Buf
	if b != nil {
		b.RefCount.Add(1)
		defer b.RefCount.Add(^uint32(0))
	}

	m.Mu.Unlock()

	if b != nil {
		q := b.Queue
		done := c.Done()

		for i, j := 0, q.Len(); i < j; i++ {
			select {
			default:
			case <-done:
				sink.Error(c.Err())
				return
			}

			Try1(sink, Next(q.At(i)), func() { sink.Error(ErrOops) })
		}
	}

	switch lastn.Kind {
	case KindError:
		sink.Error(lastn.Error)
	case KindComplete:
		sink.Complete()
	}
}
