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
	m := &multicast[T]{cap: n}
	return Subject[T]{
		Observable: NewObservable(m.Subscribe),
		Observer:   WithRuntimeFinalizer(m.Emit),
	}
}

type multicast[T any] struct {
	mu    sync.Mutex
	cap   int
	mobs  multiObserver[T]
	lastn Notification[struct{}]
	buf   *struct {
		queue.Queue[T]
		refcount atomic.Uint32
	}
}

func pnew[T any](*T) *T { return new(T) }

func (m *multicast[T]) Emit(n Notification[T]) {
	m.mu.Lock()

	if m.lastn.Kind != 0 {
		m.mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
		mobs := m.mobs.Clone()
		defer mobs.Release()

		if m.cap != 0 {
			buf := m.buf

			switch {
			case buf == nil:
				buf = pnew(buf)
				m.buf = buf
			case buf.refcount.Load() != 0:
				q := buf.Queue.Clone()
				buf = pnew(buf)
				buf.Queue = q
				m.buf = buf
			}

			q := &buf.Queue

			if q.Len() == m.cap {
				q.Pop()
			}

			q.Push(n.Value)
		}

		m.mu.Unlock()

		mobs.Emit(n)

	case KindError, KindComplete:
		var mobs multiObserver[T]

		m.mobs, mobs = mobs, m.mobs

		switch n.Kind {
		case KindError:
			m.lastn = Error[struct{}](n.Error)
		case KindComplete:
			m.lastn = Complete[struct{}]()
		}

		m.mu.Unlock()

		mobs.Emit(n)

	default: // Unknown kind.
		m.mu.Unlock()
	}
}

func (m *multicast[T]) Subscribe(c Context, o Observer[T]) {
	m.mu.Lock()

	if buf := m.buf; buf != nil {
		buf.refcount.Add(1)
		decrease := true
		defer func() {
			if decrease {
				buf.refcount.Add(^uint32(0))
			}
		}()

		m.mu.Unlock()

		q := buf.Queue
		done := c.Done()

		for i, j := 0, q.Len(); i < j; i++ {
			select {
			default:
			case <-done:
				o.Error(c.Cause())
				return
			}

			Try1(o, Next(q.At(i)), func() { o.Error(ErrOops) })
		}

		buf.refcount.Add(^uint32(0))
		decrease = false

		m.mu.Lock()
	}

	lastn := m.lastn
	if lastn.Kind == 0 {
		c, o = Serialize(c, o)

		o := o
		m.mobs.Add(&o)

		c.AfterFunc(func() {
			m.mu.Lock()
			m.mobs.Delete(&o)
			m.mu.Unlock()
			o.Error(c.Cause())
		})
	}

	m.mu.Unlock()

	switch lastn.Kind {
	case KindError:
		o.Error(lastn.Error)
	case KindComplete:
		o.Complete()
	}
}
