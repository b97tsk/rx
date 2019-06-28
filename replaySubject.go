package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/x/queue"
)

// A ReplaySubject buffers a set number of values and will emit those values
// immediately to any new subscribers in addition to emitting new values to
// existing subscribers.
type ReplaySubject struct {
	*replaySubject
}

// NewReplaySubject creates a new ReplaySubject.
func NewReplaySubject(bufferSize int, windowTime time.Duration) ReplaySubject {
	s := &replaySubject{
		BufferSize: bufferSize,
		WindowTime: windowTime,
	}
	s.Subject = Subject{
		Observable: Observable{}.Lift(s.call),
		Observer:   s.notify,
	}
	s.lock = make(chan struct{}, 1)
	s.lock <- struct{}{}
	return ReplaySubject{s}
}

type replaySubject struct {
	Subject
	lock       chan struct{}
	observers  observerList
	cws        contextWaitService
	err        error
	buffer     queue.Queue
	BufferSize int
	WindowTime time.Duration
}

type replaySubjectValue struct {
	Deadline time.Time
	Value    interface{}
}

func (s *replaySubject) trimBuffer() {
	if s.BufferSize > 0 {
		for s.buffer.Len() > s.BufferSize {
			s.buffer.PopFront()
		}
	}
	if s.WindowTime > 0 {
		now := time.Now()
		for s.buffer.Len() > 0 {
			if s.buffer.Front().(replaySubjectValue).Deadline.After(now) {
				break
			}
			s.buffer.PopFront()
		}
	}
}

func (s *replaySubject) notify(t Notification) {
	if _, ok := <-s.lock; ok {
		switch {
		case t.HasValue:
			observers, releaseRef := s.observers.AddRef()

			var deadline time.Time
			if s.WindowTime > 0 {
				deadline = time.Now().Add(s.WindowTime)
			}
			s.buffer.PushBack(replaySubjectValue{deadline, t.Value})
			s.trimBuffer()

			s.lock <- struct{}{}

			for _, sink := range observers {
				sink.Notify(t)
			}

			releaseRef()

		case t.HasError:
			observers := s.observers.Swap(nil)
			s.err = t.Value.(error)

			close(s.lock)

			for _, sink := range observers {
				sink.Notify(t)
			}

		default:
			observers := s.observers.Swap(nil)

			close(s.lock)

			for _, sink := range observers {
				sink.Notify(t)
			}
		}
	}
}

func (s *replaySubject) call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	if _, ok := <-s.lock; ok {
		ctx, cancel := context.WithCancel(ctx)

		observer := Mutex(Finally(sink, cancel))
		s.observers.Append(&observer)

		finalize := func() {
			if _, ok := <-s.lock; ok {
				s.observers.Remove(&observer)
				s.lock <- struct{}{}
			}
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = newContextWaitService()
		}

		s.trimBuffer()

		for i, j := 0, s.buffer.Len(); i < j; i++ {
			if isDone(ctx) {
				break
			}
			sink.Next(s.buffer.At(i).(replaySubjectValue).Value)
		}

		s.lock <- struct{}{}
		return ctx, cancel
	}

	if s.err != nil {
		sink.Error(s.err)
		return Done()
	}

	s.trimBuffer()

	for i, j := 0, s.buffer.Len(); i < j; i++ {
		if isDone(ctx) {
			return Done()
		}
		sink.Next(s.buffer.At(i).(replaySubjectValue).Value)
	}
	sink.Complete()
	return Done()
}
