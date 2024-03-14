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

	m := &multicastReplay[T]{n: n}

	return Subject[T]{
		Observable: NewObservable(m.subscribe),
		Observer:   NewObserver(m.emit).WithRuntimeFinalizer(),
	}
}

type multicastReplay[T any] struct {
	multicast[T]
	n int
	b struct {
		q  queue.Queue[T]
		rc *atomic.Uint32
	}
}

func (m *multicastReplay[T]) bufferForRead() (queue.Queue[T], *atomic.Uint32) {
	rc := m.b.rc

	if rc == nil {
		rc = new(atomic.Uint32)
		m.b.rc = rc
	}

	rc.Add(1)

	return m.b.q, rc
}

func (m *multicastReplay[T]) bufferForWrite() *queue.Queue[T] {
	rc := m.b.rc

	if rc != nil && rc.Load() != 0 {
		m.b.q = m.b.q.Clone()
		m.b.rc = nil
	}

	return &m.b.q
}

func (m *multicastReplay[T]) emit(n Notification[T]) {
	m.mu.Lock()

	if m.last.Kind != 0 {
		m.mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
		obs := m.obs.Clone()
		defer obs.Release()

		b := m.bufferForWrite()

		if n := m.n; n > 0 && n == b.Len() {
			b.Pop()
		}

		b.Push(n.Value)

		m.mu.Unlock()

		obs.Emit(n)

	case KindError, KindComplete:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.last = n

		if n.Kind == KindError {
			m.b.q.Init()
		}

		m.mu.Unlock()

		obs.Emit(n)

	default: // Unknown kind.
		m.mu.Unlock()
	}
}

func (m *multicastReplay[T]) subscribe(c Context, sink Observer[T]) {
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

	b, rc := m.bufferForRead()
	defer rc.Add(^uint32(0))

	m.mu.Unlock()

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

	if last.Kind != 0 {
		sink(last)
	}
}
