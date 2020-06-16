package rx

import (
	"context"
	"sync"
	"time"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
	"github.com/b97tsk/rx/internal/queue"
)

// A ReplaySubject buffers a set number of values and will emit those values
// immediately to any new subscribers in addition to emitting new values to
// existing subscribers.
type ReplaySubject struct {
	*replaySubject
}

type replaySubject struct {
	Subject
	mux        sync.Mutex
	lst        observerList
	cws        misc.ContextWaitService
	err        error
	buffer     queue.Queue
	bufferRefs *atomic.Uint32s
	bufferSize int
	windowTime time.Duration
}

type replaySubjectElement struct {
	Deadline time.Time
	Value    interface{}
}

// NewReplaySubject creates a new ReplaySubject.
func NewReplaySubject(bufferSize int, windowTime time.Duration) ReplaySubject {
	s := &replaySubject{
		bufferSize: bufferSize,
		windowTime: windowTime,
	}
	s.Subject = Subject{Create(s.subscribe), s.sink}
	return ReplaySubject{s}
}

// Exists reports if this ReplaySubject is ready to use.
func (s ReplaySubject) Exists() bool {
	return s.replaySubject != nil
}

// BufferSize returns current buffer size setting.
func (s ReplaySubject) BufferSize() int {
	return s.bufferSize
}

// WindowTime returns current window time setting.
func (s ReplaySubject) WindowTime() time.Duration {
	return s.windowTime
}

func (s *replaySubject) bufferForRead() (queue.Queue, *atomic.Uint32s) {
	refs := s.bufferRefs
	if refs == nil {
		refs = new(atomic.Uint32s)
		s.bufferRefs = refs
	}
	refs.Add(1)
	return s.buffer, refs
}

func (s *replaySubject) bufferForWrite() *queue.Queue {
	refs := s.bufferRefs
	if refs != nil && !refs.Equals(0) {
		s.buffer = s.buffer.Clone()
		s.bufferRefs = nil
	}
	return &s.buffer
}

func (s *replaySubject) trimBuffer(b *queue.Queue) {
	if s.windowTime > 0 {
		if b == nil {
			b = s.bufferForWrite()
		}
		now := time.Now()
		for b.Len() > 0 {
			if b.Front().(replaySubjectElement).Deadline.After(now) {
				break
			}
			b.Pop()
		}
	}
	if s.bufferSize > 0 {
		if b == nil {
			b = s.bufferForWrite()
		}
		for b.Len() > s.bufferSize {
			b.Pop()
		}
	}
}

func (s *replaySubject) sink(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		lst := s.lst.Clone()
		defer lst.Release()

		var deadline time.Time
		if s.windowTime > 0 {
			deadline = time.Now().Add(s.windowTime)
		}
		b := s.bufferForWrite()
		b.Push(replaySubjectElement{deadline, t.Value})
		s.trimBuffer(b)

		s.mux.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList
		s.lst.Swap(&lst)

		if t.HasError {
			s.err = t.Error
			s.buffer.Init()
		} else {
			s.err = Completed
		}

		s.mux.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (s *replaySubject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	err := s.err
	if err == nil {
		observer := Mutex(sink)
		s.lst.Append(&observer)

		finalize := func() {
			s.mux.Lock()
			s.lst.Remove(&observer)
			s.mux.Unlock()
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = misc.NewContextWaitService()
		}
	}

	s.trimBuffer(nil)

	b, refs := s.bufferForRead()
	defer refs.Sub(1)

	s.mux.Unlock()

	for i, j := 0, b.Len(); i < j; i++ {
		if ctx.Err() != nil {
			return
		}
		sink.Next(b.At(i).(replaySubjectElement).Value)
	}

	if err != nil {
		if err != Completed {
			sink.Error(err)
		} else {
			sink.Complete()
		}
	}
}
