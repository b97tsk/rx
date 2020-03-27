package rx

import (
	"context"
	"sync"
	"sync/atomic"
)

// A BehaviorSubject stores the latest value emitted to its consumers, and
// whenever a new Observer subscribes, it will immediately receive the
// "current value" from the BehaviorSubject.
type BehaviorSubject struct {
	*behaviorSubject
}

type behaviorSubject struct {
	Subject
	mux       sync.Mutex
	observers observerList
	cws       contextWaitService
	err       error
	val       atomic.Value
}

type behaviorSubjectValue struct {
	Value interface{}
}

// NewBehaviorSubject creates a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) BehaviorSubject {
	s := new(behaviorSubject)
	s.Observable = s.subscribe
	s.Observer = s.notify
	s.val.Store(behaviorSubjectValue{val})
	return BehaviorSubject{s}
}

// Exists reports if this BehaviorSubject is ready to use.
func (s BehaviorSubject) Exists() bool {
	return s.behaviorSubject != nil
}

// Value returns the latest value stored in this BehaviorSubject.
func (s BehaviorSubject) Value() interface{} {
	return s.getValue()
}

func (s *behaviorSubject) getValue() interface{} {
	return s.val.Load().(behaviorSubjectValue).Value
}

func (s *behaviorSubject) notify(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		observers, releaseRef := s.observers.AddRef()
		s.val.Store(behaviorSubjectValue{t.Value})
		s.mux.Unlock()

		for _, sink := range observers {
			sink.Notify(t)
		}

		releaseRef()

	default:
		observers := s.observers.Swap(nil)
		if t.HasError {
			s.err = t.Error
		} else {
			s.err = Complete
		}
		s.mux.Unlock()

		for _, sink := range observers {
			sink.Notify(t)
		}
	}
}

func (s *behaviorSubject) subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	s.mux.Lock()

	ctx := NewContext(parent)

	if err := s.err; err != nil {
		if err != Complete {
			sink.Error(err)
		} else {
			sink.Next(s.getValue())
			sink.Complete()
		}
		ctx.Unsubscribe(err)
	} else {
		observer := Mutex(DoAtLast(sink, ctx.AtLast))
		s.observers.Append(&observer)

		finalize := func() {
			s.mux.Lock()
			s.observers.Remove(&observer)
			s.mux.Unlock()
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = newContextWaitService()
		}

		sink.Next(s.getValue())
	}

	s.mux.Unlock()
	return ctx, ctx.Cancel
}
