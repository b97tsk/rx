package rx

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/x/misc"
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
	cws       misc.ContextWaitService
	err       error
	val       atomic.Value
}

type behaviorSubjectElement struct {
	Value interface{}
}

// NewBehaviorSubject creates a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) BehaviorSubject {
	s := new(behaviorSubject)
	s.Observable = Create(s.subscribe)
	s.Observer = s.sink
	s.val.Store(behaviorSubjectElement{val})
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
	return s.val.Load().(behaviorSubjectElement).Value
}

func (s *behaviorSubject) sink(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		observers, releaseRef := s.observers.AddRef()
		s.val.Store(behaviorSubjectElement{t.Value})
		s.mux.Unlock()

		for _, observer := range observers {
			observer.Sink(t)
		}

		releaseRef()

	default:
		observers := s.observers.Swap(nil)
		if t.HasError {
			s.err = t.Error
		} else {
			s.err = Completed
		}
		s.mux.Unlock()

		for _, observer := range observers {
			observer.Sink(t)
		}
	}
}

func (s *behaviorSubject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	if err := s.err; err != nil {
		if err != Completed {
			sink.Error(err)
		} else {
			sink.Next(s.getValue())
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

		sink.Next(s.getValue())
	}

	s.mux.Unlock()
}
