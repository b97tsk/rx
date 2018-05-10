package rx

import (
	"context"
)

// A BehaviorSubject stores the latest value emitted to its consumers, and
// whenever a new Observer subscribes, it will immediately receive the
// "current value" from the BehaviorSubject.
type BehaviorSubject struct {
	Subject
	value interface{}
}

// BehaviorSubject implements SubjectLike.
var _ SubjectLike = (*BehaviorSubject)(nil)

// Next stores the latest value and emits it to the consumers of this
// BehaviorSubject.
func (s *BehaviorSubject) Next(val interface{}) {
	if s.try.Lock() {
		defer s.try.Unlock()
		s.value = val
		for _, ob := range s.observers {
			ob.Next(val)
		}
	}
}

// Value returns the latest value stored in this BehaviorSubject.
func (s *BehaviorSubject) Value() interface{} {
	if s.try.Lock() {
		val := s.value
		s.try.Unlock()
		return val
	}
	return s.value
}

// Subscribe adds a consumer to this BehaviorSubject.
func (s *BehaviorSubject) Subscribe(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	if s.try.Lock() {
		defer s.try.Unlock()

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
func NewBehaviorSubject(val interface{}) *BehaviorSubject {
	s := &BehaviorSubject{value: val}
	s.Op = OperatorFunc(s.Subscribe)
	return s
}
