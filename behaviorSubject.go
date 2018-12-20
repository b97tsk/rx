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
	s.val.Store(behaviorSubjectValue{val})
	return BehaviorSubject{s}
}

// Value returns the latest value stored in this BehaviorSubject.
func (s BehaviorSubject) Value() interface{} {
	return s.getValue()
}

type behaviorSubject struct {
	Subject
	try       cancellableLocker
	observers []*Observer
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
	if s.try.Lock() {
		switch {
		case t.HasValue:
			s.val.Store(behaviorSubjectValue{t.Value})

			for _, sink := range s.observers {
				sink.Notify(t)
			}

			s.try.Unlock()

		case t.HasError:
			observers := s.observers
			s.observers = nil
			s.err = t.Value.(error)

			s.try.CancelAndUnlock()

			for _, sink := range observers {
				sink.Notify(t)
			}

		default:
			observers := s.observers
			s.observers = nil

			s.try.CancelAndUnlock()

			for _, sink := range observers {
				sink.Notify(t)
			}
		}
	}
}

func (s *behaviorSubject) call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	if s.try.Lock() {
		ctx, cancel := context.WithCancel(ctx)

		observer := Finally(sink, cancel)
		s.observers = append(s.observers, &observer)

		go func() {
			<-ctx.Done()
			if s.try.Lock() {
				for i, sink := range s.observers {
					if sink == &observer {
						copy(s.observers[i:], s.observers[i+1:])
						n := len(s.observers)
						s.observers[n-1] = nil
						s.observers = s.observers[:n-1]
						break
					}
				}
				s.try.Unlock()
			}
		}()

		sink.Next(s.getValue())

		s.try.Unlock()
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
