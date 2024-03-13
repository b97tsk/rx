package rx

import (
	"sync/atomic"
	"time"

	"github.com/b97tsk/rx/internal/queue"
)

// ReplayConfig carries options for MulticastReplay.
type ReplayConfig struct {
	BufferSize int
	WindowTime time.Duration
}

// MulticastReplay returns a Subject that keeps track of a certain number
// and/or a window time of values it receive. Each subscriber will then
// receive all tracked values as well as future values.
func MulticastReplay[T any](opts *ReplayConfig) Subject[T] {
	m := new(multicastReplay[T])

	if opts != nil {
		m.ReplayConfig = *opts
	}

	return Subject[T]{
		Observable: NewObservable(m.subscribe),
		Observer:   NewObserver(m.emit).WithRuntimeFinalizer(),
	}
}

type multicastReplay[T any] struct {
	multicast[T]
	ReplayConfig
	b struct {
		q  queue.Queue[Pair[time.Time, T]]
		rc *atomic.Uint32
	}
}

func (m *multicastReplay[T]) bufferForRead() (queue.Queue[Pair[time.Time, T]], *atomic.Uint32) {
	rc := m.b.rc

	if rc == nil {
		rc = new(atomic.Uint32)
		m.b.rc = rc
	}

	rc.Add(1)

	return m.b.q, rc
}

func (m *multicastReplay[T]) bufferForWrite() *queue.Queue[Pair[time.Time, T]] {
	rc := m.b.rc

	if rc != nil && rc.Load() != 0 {
		m.b.q = m.b.q.Clone()
		m.b.rc = nil
	}

	return &m.b.q
}

func (m *multicastReplay[T]) trimBuffer(b *queue.Queue[Pair[time.Time, T]]) {
	if m.WindowTime > 0 {
		if b == nil {
			b = m.bufferForWrite()
		}

		now := time.Now()

		for b.Len() != 0 {
			if b.Front().Key.After(now) {
				break
			}

			b.Pop()
		}
	}

	if n := m.BufferSize; n > 0 {
		if b == nil {
			b = m.bufferForWrite()
		}

		for b.Len() > n {
			b.Pop()
		}
	}
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

		var deadline time.Time

		if windowTime := m.WindowTime; windowTime > 0 {
			deadline = time.Now().Add(windowTime)
		}

		b := m.bufferForWrite()
		b.Push(NewPair(deadline, n.Value))
		m.trimBuffer(b)

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

	m.trimBuffer(nil)

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

		Try1(sink, Next(b.At(i).Value), func() { sink.Error(ErrOops) })
	}

	if last.Kind != 0 {
		sink(last)
	}
}
