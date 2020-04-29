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
	observers  observerList
	cws        misc.ContextWaitService
	err        error
	buffer     queue.Queue
	BufferSize int
	WindowTime time.Duration
}

type replaySubjectValue struct {
	Deadline time.Time
	Value    interface{}
}

// NewReplaySubject creates a new ReplaySubject.
func NewReplaySubject(bufferSize int, windowTime time.Duration) ReplaySubject {
	s := &replaySubject{
		BufferSize: bufferSize,
		WindowTime: windowTime,
	}
	s.Observable = Create(s.subscribe)
	s.Observer = s.notify
	return ReplaySubject{s}
}

// Exists reports if this ReplaySubject is ready to use.
func (s ReplaySubject) Exists() bool {
	return s.replaySubject != nil
}

func (s *replaySubject) trimBuffer() {
	if s.WindowTime > 0 {
		now := time.Now()
		for s.buffer.Len() > 0 {
			if s.buffer.Front().(replaySubjectValue).Deadline.After(now) {
				break
			}
			s.buffer.PopFront()
		}
	}
	if s.BufferSize > 0 {
		for s.buffer.Len() > s.BufferSize {
			s.buffer.PopFront()
		}
	}
}

func (s *replaySubject) notify(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		observers, releaseRef := s.observers.AddRef()

		var deadline time.Time
		if s.WindowTime > 0 {
			deadline = time.Now().Add(s.WindowTime)
		}
		s.buffer.PushBack(replaySubjectValue{deadline, t.Value})
		s.trimBuffer()

		s.mux.Unlock()

		for _, observer := range observers {
			observer.Sink(t)
		}

		releaseRef()

	default:
		observers := s.observers.Swap(nil)
		if t.HasError {
			s.err = t.Error
			s.buffer.Init()
		} else {
			s.err = Complete
		}
		s.mux.Unlock()

		for _, observer := range observers {
			observer.Sink(t)
		}
	}
}

func (s *replaySubject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	if err := s.err; err != nil {
		if err != Complete {
			sink.Error(err)
		} else {
			s.trimBuffer()
			for i, j := 0, s.buffer.Len(); i < j; i++ {
				if ctx.Err() != nil {
					s.mux.Unlock()
					return
				}
				sink.Next(s.buffer.At(i).(replaySubjectValue).Value)
			}
			sink.Complete()
		}
	} else {
		observer := Mutex(sink)
		s.observers.Append(&observer)

		finalize := func() {
			s.mux.Lock()
			s.observers.Remove(&observer)
			s.mux.Unlock()
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = misc.NewContextWaitService()
		}

		s.trimBuffer()

		for i, j := 0, s.buffer.Len(); i < j; i++ {
			if ctx.Err() != nil {
				break
			}
			sink.Next(s.buffer.At(i).(replaySubjectValue).Value)
		}
	}

	s.mux.Unlock()
}
