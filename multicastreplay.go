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
// all Observers that subscribed to it, all of which will receive
// buffered emissions and new emissions from Subject's Observer part.
//
// Subjects are subject to memory leaks.
// After finishing using a Subject, you should call either its Error method
// or its Complete method to avoid that.
// If you can guarantee that every subscription to a Subject is canceled
// sooner or later, then you are fine.
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
			m.emit(n)

			if !n.HasValue {
				runtime.SetFinalizer(&m, nil)
			}
		},
	}
}

func multicastReplayFinalizer[T any](m **multicastReplay[T]) {
	go (*m).emit(Error[T](ErrFinalized))
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

func (m *multicastReplay[T]) emit(n Notification[T]) {
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
			observer.Emit(n)
		}

	default:
		var obs multiObserver[T]
		obs, m.obs = m.obs, obs

		m.err = errComplete
		if n.HasError {
			m.err = errOrErrNil(n.Error)
			m.b.q.Init()
		}

		m.mu.Unlock()

		for _, observer := range obs.Observers {
			observer.Emit(n)
		}
	}
}

func (m *multicastReplay[T]) subscribe(ctx context.Context, sink Observer[T]) {
	m.mu.Lock()

	err := m.err
	if err == nil {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.OnLastNotification(cancel).Serialized()

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
		select {
		default:
		case <-done:
			sink.Error(ctx.Err())
			return
		}

		sink.Next(b.At(i).Value)
	}

	if err != nil {
		if err == errComplete {
			sink.Complete()
			return
		}

		sink.Error(cleanErrNil(err))
	}
}
