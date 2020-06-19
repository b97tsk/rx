package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
	"github.com/b97tsk/rx/internal/queue"
)

// A ReplaySubject buffers a set number of values and will emit those values
// immediately to any new subscribers in addition to emitting new values to
// existing subscribers.
type ReplaySubject struct {
	Subject

	buffer     queue.Queue
	bufferRefs *atomic.Uint32s
	bufferSize atomic.Int64s
	windowTime atomic.Int64s
}

type replaySubjectElement struct {
	Deadline time.Time
	Value    interface{}
}

// NewReplaySubject creates a new ReplaySubject.
func NewReplaySubject(bufferSize int) *ReplaySubject {
	s := new(ReplaySubject)
	s.Double = Double{Create(s.subscribe), s.sink}
	s.bufferSize = atomic.Int64(int64(bufferSize))
	return s
}

// BufferSize gets current buffer size value.
func (s *ReplaySubject) BufferSize() int {
	return int(s.bufferSize.Load())
}

// WindowTime gets current window time value.
func (s *ReplaySubject) WindowTime() time.Duration {
	return time.Duration(s.windowTime.Load())
}

// SetBufferSize sets buffer size to a specified value.
func (s *ReplaySubject) SetBufferSize(bufferSize int) {
	s.bufferSize.Store(int64(bufferSize))
}

// SetWindowTime sets window time to a specified value.
func (s *ReplaySubject) SetWindowTime(windowTime time.Duration) {
	s.windowTime.Store(int64(windowTime))
}

func (s *ReplaySubject) bufferForRead() (queue.Queue, *atomic.Uint32s) {
	refs := s.bufferRefs
	if refs == nil {
		refs = new(atomic.Uint32s)
		s.bufferRefs = refs
	}
	refs.Add(1)
	return s.buffer, refs
}

func (s *ReplaySubject) bufferForWrite() *queue.Queue {
	refs := s.bufferRefs
	if refs != nil && !refs.Equals(0) {
		s.buffer = s.buffer.Clone()
		s.bufferRefs = nil
	}
	return &s.buffer
}

func (s *ReplaySubject) trimBuffer(b *queue.Queue) {
	if s.WindowTime() > 0 {
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
	if bufferSize := s.BufferSize(); bufferSize > 0 {
		if b == nil {
			b = s.bufferForWrite()
		}
		for b.Len() > bufferSize {
			b.Pop()
		}
	}
}

func (s *ReplaySubject) sink(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		lst := s.lst.Clone()
		defer lst.Release()

		var deadline time.Time
		if windowTime := s.WindowTime(); windowTime > 0 {
			deadline = time.Now().Add(windowTime)
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

func (s *ReplaySubject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	err := s.err
	if err == nil {
		observer := MutexContext(ctx, sink)
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
