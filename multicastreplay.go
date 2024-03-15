package rx

import (
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
)

// MulticastReplayAll returns a Subject that keeps track of every value
// it receives. Each subscriber will then receive all tracked values as well
// as future values.
func MulticastReplayAll[T any]() Subject[T] {
	return MulticastReplay[T](-1)
}

// MulticastReplay returns a Subject that keeps track of a certain number of
// values it receive. Each subscriber will then receive all tracked values as
// well as future values.
//
// If n < 0, MulticastReplay keeps track of every value it receives;
// if n == 0, MulticastReplay returns Multicast().
func MulticastReplay[T any](n int) Subject[T] {
	if n == 0 {
		return Multicast[T]()
	}

	m := &multicastReplay[T]{Cap: n}

	return Subject[T]{
		Observable: NewObservable(m.subscribe),
		Observer:   NewObserver(m.emit).WithRuntimeFinalizer(),
	}
}

type multicastReplay[T any] struct {
	multicast[T]
	Cap int
	Buf struct {
		Queue    queue.Queue[T]
		RefCount *atomic.Uint32
	}
}

func (m *multicastReplay[T]) bufferForRead() (queue.Queue[T], *atomic.Uint32) {
	rc := m.Buf.RefCount

	if rc == nil {
		rc = new(atomic.Uint32)
		m.Buf.RefCount = rc
	}

	rc.Add(1)

	return m.Buf.Queue, rc
}

func (m *multicastReplay[T]) bufferForWrite() *queue.Queue[T] {
	rc := m.Buf.RefCount

	if rc != nil && rc.Load() != 0 {
		m.Buf.Queue = m.Buf.Queue.Clone()
		m.Buf.RefCount = nil
	}

	return &m.Buf.Queue
}

func (m *multicastReplay[T]) emit(n Notification[T]) {
	m.Mu.Lock()

	if m.LastN.Kind != 0 {
		m.Mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
		mobs := m.Mobs.Clone()
		defer mobs.Release()

		b := m.bufferForWrite()

		if n := m.Cap; n > 0 && n == b.Len() {
			b.Pop()
		}

		b.Push(n.Value)

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

		if n.Kind == KindError {
			m.Buf.Queue.Init()
		}

		m.Mu.Unlock()

		mobs.Emit(n)

	default: // Unknown kind.
		m.Mu.Unlock()
	}
}

func (m *multicastReplay[T]) subscribe(c Context, sink Observer[T]) {
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

	b, rc := m.bufferForRead()
	defer rc.Add(^uint32(0))

	m.Mu.Unlock()

	done := c.Done()

	for i, j := 0, b.Len(); i < j; i++ {
		select {
		default:
		case <-done:
			sink.Error(c.Err())
			return
		}

		Try1(sink, Next(b.At(i)), func() { sink.Error(ErrOops) })
	}

	switch lastn.Kind {
	case KindError:
		sink.Error(lastn.Error)
	case KindComplete:
		sink.Complete()
	}
}
