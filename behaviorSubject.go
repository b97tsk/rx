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
	try       cancellableLocker
	observers []*Observer
	val       atomic.Value
	err       error
}

type behaviorSubjectValue struct {
	Value interface{}
}

// Value returns the latest value stored in this BehaviorSubject.
func (s *BehaviorSubject) Value() interface{} {
	return s.val.Load().(behaviorSubjectValue).Value
}

func (s *BehaviorSubject) notify(t Notification) {
	if s.try.Lock() {
		switch {
		case t.HasValue:
			s.val.Store(behaviorSubjectValue{t.Value})

			defer s.try.Unlock()

			for _, sink := range s.observers {
				sink.Notify(t)
			}

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

func (s *BehaviorSubject) call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	if s.try.Lock() {
		defer s.try.Unlock()

		ctx, cancel := context.WithCancel(ctx)

		observer := Finally(sink, cancel)
		s.observers = append(s.observers, &observer)

		go func() {
			<-ctx.Done()
			if s.try.Lock() {
				for i, sink := range s.observers {
					if sink == &observer {
						copy(s.observers[i:], s.observers[i+1:])
						s.observers[len(s.observers)-1] = nil
						s.observers = s.observers[:len(s.observers)-1]
						break
					}
				}
				s.try.Unlock()
			}
		}()

		sink.Next(s.Value())
		return ctx, cancel
	}

	if s.err != nil {
		sink.Error(s.err)
	} else {
		sink.Next(s.Value())
		sink.Complete()
	}

	return canceledCtx, nothingToDo
}

// NewBehaviorSubject returns a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) *BehaviorSubject {
	s := new(BehaviorSubject)
	s.Observer = s.notify
	s.Observable = s.Observable.Lift(s.call)
	s.val.Store(behaviorSubjectValue{val})
	return s
}
