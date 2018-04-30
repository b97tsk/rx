package rx

import (
	"context"
)

// A BehaviorSubject stores the latest value emitted to its consumers, and
// whenever a new Observer subscribes, it will immediately receive the
// "current value" from the BehaviorSubject.
type BehaviorSubject interface {
	Subject
	Value() interface{}
}

type behaviorSubjectImplement struct {
	subjectImplement
	value interface{}
}

func (s *behaviorSubjectImplement) Next(val interface{}) {
	if s.try.Lock() {
		s.value = val
		for _, ob := range s.observers {
			ob.Next(val)
		}
		s.try.Unlock()
	}
}

func (s *behaviorSubjectImplement) Value() interface{} {
	if s.try.Lock() {
		val := s.value
		s.try.Unlock()
		return val
	}
	return s.value
}

func (s *behaviorSubjectImplement) AsObservable() Observable {
	return Observable{OperatorFunc(s.Subscribe)}
}

func (s *behaviorSubjectImplement) Subscribe(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	if s.try.Lock() {
		ctx, cancel := context.WithCancel(ctx)

		observer := withFinalizer(ob, cancel)
		s.observers = append(s.observers, &observer)

		go func() {
			<-ctx.Done()
			if s.try.Lock() {
				for i, ob := range s.observers {
					if ob == &observer {
						copy(s.observers[i:], s.observers[i+1:])
						s.observers[len(s.observers)-1] = nil
						s.observers = s.observers[:len(s.observers)-1]
						break
					}
				}
				s.try.Unlock()
			}
		}()

		ob.Next(s.value)
		s.try.Unlock()

		return ctx, cancel
	}

	if s.hasError {
		ob.Error(s.errValue)
	} else {
		ob.Next(s.value)
		ob.Complete()
	}

	return canceledCtx, noopFunc
}

// NewBehaviorSubject returns a new BehaviorSubject.
func NewBehaviorSubject(val interface{}) BehaviorSubject {
	return &behaviorSubjectImplement{value: val}
}
