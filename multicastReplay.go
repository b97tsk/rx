package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/atomic"
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
	s := &multicastReplay{}

	if opts != nil {
		s.ReplayOptions = *opts
	}

	return Subject{
		Observable: s.subscribe,
		Observer:   s.sink,
	}
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

func (s *multicastReplay) bufferForRead() (queue.Queue, *atomic.Uint32) {
	refs := s.bufferRefs

	if refs == nil {
		refs = new(atomic.Uint32)
		s.bufferRefs = refs
	}

	refs.Add(1)

	return s.buffer, refs
}

func (s *multicastReplay) bufferForWrite() *queue.Queue {
	refs := s.bufferRefs

	if refs != nil && !refs.Equals(0) {
		s.buffer = s.buffer.Clone()
		s.bufferRefs = nil
	}

	return &s.buffer
}

func (s *multicastReplay) trimBuffer(b *queue.Queue) {
	if s.WindowTime > 0 {
		if b == nil {
			b = s.bufferForWrite()
		}

		now := time.Now()

		for b.Len() > 0 {
			if b.Front().(multicastReplayElement).Deadline.After(now) {
				break
			}

			b.Pop()
		}
	}

	if bufferSize := s.BufferSize; bufferSize > 0 {
		if b == nil {
			b = s.bufferForWrite()
		}

		for b.Len() > bufferSize {
			b.Pop()
		}
	}
}

func (s *multicastReplay) sink(t Notification) {
	s.mu.Lock()

	switch {
	case s.err != nil:
		s.mu.Unlock()

	case t.HasValue:
		lst := s.lst.Clone()
		defer lst.Release()

		var deadline time.Time

		if windowTime := s.WindowTime; windowTime > 0 {
			deadline = time.Now().Add(windowTime)
		}

		b := s.bufferForWrite()
		b.Push(multicastReplayElement{t.Value, deadline})
		s.trimBuffer(b)

		s.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList

		s.lst.Swap(&lst)

		s.err = errCompleted

		if t.HasError {
			s.err = t.Error

			if s.err == nil {
				s.err = errNil
			}

			s.buffer.Init()
		}

		s.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (s *multicastReplay) subscribe(ctx context.Context, sink Observer) {
	s.mu.Lock()

	err := s.err
	if err == nil {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		sink = sink.WithCancel(cancel).MutexContext(ctx)

		observer := sink
		s.lst.Append(&observer)
		s.cws.Submit(ctx, func() {
			s.mu.Lock()
			s.lst.Remove(&observer)
			s.mu.Unlock()
		})
	}

	s.trimBuffer(nil)

	b, refs := s.bufferForRead()
	defer refs.Sub(1)

	s.mu.Unlock()

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
