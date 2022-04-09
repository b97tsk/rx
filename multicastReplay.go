package rx

import (
	"context"
	"runtime"
	"time"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/ctxwatch"
	"github.com/b97tsk/rx/internal/queue"
)

// ReplayOptions are the options for MulticastReplay.
type ReplayOptions struct {
	BufferSize int
	WindowTime time.Duration
}

// MulticastReplay returns a Subject whose Observable part takes care of all
// Observers that subscribes to it, which will receive buffered emissions and
// new emissions from Subject's Observer part.
func MulticastReplay(opts *ReplayOptions) Subject {
	m := &multicastReplay{}

	if opts != nil {
		m.ReplayOptions = *opts
	}

	runtime.SetFinalizer(&m, multicastReplayFinalizer)

	return Subject{
		Observable: func(ctx context.Context, sink Observer) {
			m.subscribe(ctx, sink)
		},
		Observer: func(t Notification) {
			m.sink(t)
		},
	}
}

func multicastReplayFinalizer(m **multicastReplay) {
	go Observer((*m).sink).Error(errFinalized)
}

// MulticastReplayFactory returns a SubjectFactory that wraps calls to
// MulticastReplay.
func MulticastReplayFactory(opts *ReplayOptions) SubjectFactory {
	return func() Subject { return MulticastReplay(opts) }
}

type multicastReplay struct {
	ReplayOptions
	multicast
	buffer     queue.Queue
	bufferRefs *atomic.Uint32
}

type multicastReplayElement struct {
	Value    interface{}
	Deadline time.Time
}

func (m *multicastReplay) bufferForRead() (queue.Queue, *atomic.Uint32) {
	refs := m.bufferRefs

	if refs == nil {
		refs = new(atomic.Uint32)
		m.bufferRefs = refs
	}

	refs.Add(1)

	return m.buffer, refs
}

func (m *multicastReplay) bufferForWrite() *queue.Queue {
	refs := m.bufferRefs

	if refs != nil && !refs.Equals(0) {
		m.buffer = m.buffer.Clone()
		m.bufferRefs = nil
	}

	return &m.buffer
}

func (m *multicastReplay) trimBuffer(b *queue.Queue) {
	if m.WindowTime > 0 {
		if b == nil {
			b = m.bufferForWrite()
		}

		now := time.Now()

		for b.Len() > 0 {
			if b.Front().(multicastReplayElement).Deadline.After(now) {
				break
			}

			b.Pop()
		}
	}

	if bufferSize := m.BufferSize; bufferSize > 0 {
		if b == nil {
			b = m.bufferForWrite()
		}

		for b.Len() > bufferSize {
			b.Pop()
		}
	}
}

func (m *multicastReplay) sink(t Notification) {
	m.mu.Lock()

	switch {
	case m.err != nil:
		m.mu.Unlock()

	case t.HasValue:
		lst := m.lst.Clone()
		defer lst.Release()

		var deadline time.Time

		if windowTime := m.WindowTime; windowTime > 0 {
			deadline = time.Now().Add(windowTime)
		}

		b := m.bufferForWrite()
		b.Push(multicastReplayElement{t.Value, deadline})
		m.trimBuffer(b)

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

			m.buffer.Init()
		}

		m.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (m *multicastReplay) subscribe(ctx context.Context, sink Observer) {
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

	m.trimBuffer(nil)

	b, refs := m.bufferForRead()
	defer refs.Sub(1)

	m.mu.Unlock()

	for i, j := 0, b.Len(); i < j; i++ {
		if ctx.Err() != nil {
			return
		}

		sink.Next(b.At(i).(multicastReplayElement).Value)
	}

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
