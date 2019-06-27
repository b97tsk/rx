package rx

import (
	"context"
	"sync/atomic"
)

// A BehaviorSubject stores the latest value emitted to its consumers, and
// whenever a new Observer subscribes, it will immediately receive the
// "current value" from the BehaviorSubject.
type BehaviorSubject struct {
	*behaviorSubject
}

// NewBehaviorSubject creates a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) BehaviorSubject {
	s := new(behaviorSubject)
	s.Subject = Subject{
		Observable: Observable{}.Lift(s.call),
		Observer:   s.notify,
	}
	s.lock = make(chan struct{}, 1)
	s.lock <- struct{}{}
	s.val.Store(behaviorSubjectValue{val})
	return BehaviorSubject{s}
}

// Value returns the latest value stored in this BehaviorSubject.
func (s BehaviorSubject) Value() interface{} {
	return s.getValue()
}

type behaviorSubject struct {
	Subject
	lock      chan struct{}
	observers observerList
	err       error
	val       atomic.Value
}

type behaviorSubjectValue struct {
	Value interface{}
}

func (s *behaviorSubject) getValue() interface{} {
	return s.val.Load().(behaviorSubjectValue).Value
}

func (s *behaviorSubject) notify(t Notification) {
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

func (s *behaviorSubject) call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	if _, ok := <-s.lock; ok {
		ctx, cancel := context.WithCancel(ctx)

		observer := Mutex(Finally(sink, cancel))
		s.observers.Append(&observer)

		go func() {
			<-ctx.Done()
			if _, ok := <-s.lock; ok {
				s.observers.Remove(&observer)
				s.lock <- struct{}{}
			}
		}()

		sink.Next(s.getValue())

		s.lock <- struct{}{}
		return ctx, cancel
	}

	if s.err != nil {
		sink.Error(s.err)
	} else {
		sink.Next(s.getValue())
		sink.Complete()
	}

	return Done()
}
