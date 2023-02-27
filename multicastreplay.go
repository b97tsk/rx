package rx

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/b97tsk/rx/internal/ctxwatch"
	"github.com/b97tsk/rx/internal/queue"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// ReplayConfig carries options for MulticastReplay.
type ReplayConfig struct {
	BufferSize int
	WindowTime time.Duration
}

// MulticastReplay returns a Subject whose Observable part takes care of
// all Observers that subscribed to it, all of which will receive buffered
// emissions and new emissions from Subject's Observer part.
//
// Subjects are subject to memory leaks. Once you have done using a Subject,
// you should at least call its Error or Complete method once to avoid that.
// If you can guarantee that every subscription to the Subject is canceled,
// then you are fine.
//
// Internally, a runtime finalizer is set to call Error(ErrFinalized), but
// it is not guaranteed to work.
func MulticastReplay[T any](opts *ReplayConfig) Subject[T] {
	m := new(multicastReplay[T])

	if opts != nil {
		m.ReplayConfig = *opts
	}

	runtime.SetFinalizer(&m, multicastReplayFinalizer[T])

	return Subject[T]{
		Observable: func(ctx context.Context, sink Observer[T]) {
			m.subscribe(ctx, sink)
		},
		Observer: func(n Notification[T]) {
			m.sink(n)
		},
	}
}

func multicastReplayFinalizer[T any](m **multicastReplay[T]) {
	go (*m).sink(Error[T](ErrFinalized))
}

type multicastReplay[T any] struct {
	multicast[T]
	ReplayConfig
	buffer   queue.Queue[Pair[time.Time, T]]
	bufferRc *atomic.Uint32
}

func (m *multicastReplay[T]) bufferForRead() (queue.Queue[Pair[time.Time, T]], *atomic.Uint32) {
	rc := m.bufferRc

	if rc == nil {
		rc = new(atomic.Uint32)
		m.bufferRc = rc
	}

	rc.Add(1)

	return m.buffer, rc
}

func (m *multicastReplay[T]) bufferForWrite() *queue.Queue[Pair[time.Time, T]] {
	rc := m.bufferRc

	if rc != nil && rc.Load() != 0 {
		m.buffer = m.buffer.Clone()
		m.bufferRc = nil
	}

	return &m.buffer
}

func (m *multicastReplay[T]) trimBuffer(b *queue.Queue[Pair[time.Time, T]]) {
	if m.WindowTime > 0 {
		if b == nil {
			b = m.bufferForWrite()
		}

		now := time.Now()

		for b.Len() > 0 {
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

func (m *multicastReplay[T]) sink(n Notification[T]) {
	m.mu.Lock()

	switch {
	case m.err != nil:
		m.mu.Unlock()

	case n.HasValue:
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

		for _, observer := range obs.Observers {
			observer.Sink(n)
		}

	default:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.err = errCompleted
		if n.HasError {
			m.err = errOrErrNil(n.Error)
			m.buffer.Init()
		}

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Sink(n)
		}
	}
}

func (m *multicastReplay[T]) subscribe(ctx context.Context, sink Observer[T]) {
	m.mu.Lock()

	err := m.err
	if err == nil {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.WithCancel(cancel).WithMutex()

		observer := sink
		m.obs.Add(&observer)

		wg := waitgroup.Get(ctx)
		if wg != nil {
			wg.Add(1)
		}

		ctxwatch.Add(ctx, func() {
			if wg != nil {
				defer wg.Done()
			}

			m.mu.Lock()
			m.obs.Delete(&observer)
			m.mu.Unlock()

			observer.Error(ctx.Err())
		})
	}

	m.trimBuffer(nil)

	b, rc := m.bufferForRead()
	defer rc.Add(^uint32(0))

	m.mu.Unlock()

	done := ctx.Done()

	for i, j := 0, b.Len(); i < j; i++ {
		if err := getErrWithDoneChan(ctx, done); err != nil {
			sink.Error(err)
			return
		}

		sink.Next(b.At(i).Value)
	}

	if err != nil {
		if err == errCompleted {
			sink.Complete()
			return
		}

		sink.Error(cleanErrNil(err))
	}
}
