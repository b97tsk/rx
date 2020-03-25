package rx

import (
	"context"
	"sync/atomic"
)

// A BehaviorSubject stores the latest value emitted to its consumers, and
// whenever a new Observer subscribes, it will immediately receive the
// "current value" from the BehaviorSubject.
type BehaviorSubject struct {
	Subject
	lock      chan struct{}
	observers observerList
	cws       contextWaitService
	err       error
	val       atomic.Value
}

// NewBehaviorSubject creates a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) *BehaviorSubject {
	s := new(BehaviorSubject)
	s.Observable = s.subscribe
	s.Observer = s.notify
	s.lock = make(chan struct{}, 1)
	s.lock <- struct{}{}
	s.val.Store(behaviorSubjectValue{val})
	return s
}

type behaviorSubjectValue struct {
	Value interface{}
}

// Value returns the latest value stored in this BehaviorSubject.
func (s *BehaviorSubject) Value() interface{} {
	return s.val.Load().(behaviorSubjectValue).Value
}

func (s *BehaviorSubject) notify(t Notification) {
	if _, ok := <-s.lock; ok {
		switch {
		case t.HasValue:
			observers, releaseRef := s.observers.AddRef()

			s.val.Store(behaviorSubjectValue{t.Value})

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

func (s *BehaviorSubject) subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
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

		sink.Next(s.Value())

		s.lock <- struct{}{}
		return ctx, ctx.Cancel
	}

	if s.err != nil {
		sink.Error(s.err)
		ctx.Unsubscribe(s.err)
	} else {
		sink.Next(s.Value())
		sink.Complete()
		ctx.Unsubscribe(Complete)
	}

	return ctx, ctx.Cancel
}
