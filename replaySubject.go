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
	Subject
	BufferSize int
	WindowTime time.Duration
	lock       chan struct{}
	observers  observerList
	cws        contextWaitService
	err        error
	buffer     queue.Queue
}

// NewReplaySubject creates a new ReplaySubject.
func NewReplaySubject(bufferSize int, windowTime time.Duration) *ReplaySubject {
	s := &ReplaySubject{
		BufferSize: bufferSize,
		WindowTime: windowTime,
	}
	s.Observable = s.subscribe
	s.Observer = s.notify
	s.lock = make(chan struct{}, 1)
	s.lock <- struct{}{}
	return s
}

type replaySubjectValue struct {
	Deadline time.Time
	Value    interface{}
}

func (s *ReplaySubject) trimBuffer() {
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

func (s *ReplaySubject) notify(t Notification) {
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
			s.err = t.Error

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

func (s *ReplaySubject) subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	if _, ok := <-s.lock; ok {
		observer := Mutex(DoAtLast(sink, ctx.AtLast))
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
			if ctx.Err() != nil {
				break
			}
			sink.Next(s.buffer.At(i).(replaySubjectValue).Value)
		}

		s.lock <- struct{}{}
		return ctx, ctx.Cancel
	}

	if s.err != nil {
		sink.Error(s.err)
		ctx.Unsubscribe(s.err)
		return ctx, ctx.Cancel
	}

	s.trimBuffer()

	for i, j := 0, s.buffer.Len(); i < j; i++ {
		if ctx.Err() != nil {
			return ctx, ctx.Cancel
		}
		sink.Next(s.buffer.At(i).(replaySubjectValue).Value)
	}
	sink.Complete()
	ctx.Unsubscribe(Complete)
	return ctx, ctx.Cancel
}
