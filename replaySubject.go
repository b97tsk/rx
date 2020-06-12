package rx

import (
	"context"
	"sync"
	"time"

	"github.com/b97tsk/rx/x/misc"
	"github.com/b97tsk/rx/x/queue"
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

func (s *replaySubject) trimBuffer() {
	if s.windowTime > 0 {
		now := time.Now()
		for s.buffer.Len() > 0 {
			if s.buffer.Front().(replaySubjectElement).Deadline.After(now) {
				break
			}
			s.buffer.PopFront()
		}
	}
	if s.bufferSize > 0 {
		for s.buffer.Len() > s.bufferSize {
			s.buffer.PopFront()
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
		s.buffer.PushBack(replaySubjectElement{deadline, t.Value})
		s.trimBuffer()

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

	s.trimBuffer()

	for i, j := 0, s.buffer.Len(); i < j; i++ {
		if ctx.Err() != nil {
			s.mux.Unlock()
			return
		}
		sink.Next(s.buffer.At(i).(replaySubjectElement).Value)
	}

	s.mux.Unlock()

	if err != nil {
		if err != Completed {
			sink.Error(err)
		} else {
			sink.Complete()
		}
	}
}
