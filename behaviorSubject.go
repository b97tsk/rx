package rx

import (
	"context"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/misc"
)

// A BehaviorSubject stores the latest value emitted to its consumers, and
// whenever a new Observer subscribes, it will immediately receive the
// "current value" from the BehaviorSubject.
type BehaviorSubject struct {
	Subject

	val atomic.Value
}

type behaviorSubjectElement struct {
	Value interface{}
}

// NewBehaviorSubject creates a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) *BehaviorSubject {
	s := new(BehaviorSubject)
	s.Double = Double{s.subscribe, s.sink}
	s.val.Store(behaviorSubjectElement{val})
	return s
}

// Value returns the latest value stored in this BehaviorSubject.
func (s *BehaviorSubject) Value() interface{} {
	return s.val.Load().(behaviorSubjectElement).Value
}

func (s *BehaviorSubject) sink(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		lst := s.lst.Clone()
		defer lst.Release()

		s.val.Store(behaviorSubjectElement{t.Value})

		s.mux.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList
		s.lst.Swap(&lst)

		if t.HasError {
			s.err = t.Error
		} else {
			s.err = errCompleted
		}

		s.mux.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (s *BehaviorSubject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	err := s.err
	if err == nil {
		ctx, cancel := context.WithCancel(ctx)
		observer := sink.WithCancel(cancel).MutexContext(ctx)
		s.lst.Append(&observer)

		finalize := func() {
			s.mux.Lock()
			s.lst.Remove(&observer)
			s.mux.Unlock()
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = misc.NewContextWaitService()
		}

		sink.Next(s.Value())
	}

	s.mux.Unlock()

	if err != nil {
		if err != errCompleted {
			sink.Error(err)
		} else {
			sink.Next(s.Value())
			sink.Complete()
		}
	}
}
