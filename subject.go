package rx

import (
	"context"
)

// SubjectLike is the interface represents a Subject.
type SubjectLike interface {
	Observer
	Subscribe(context.Context, Observer) (context.Context, context.CancelFunc)
}

// Subject is a special type of Observable that allows values to be multicasted
// to many Observers.
type Subject struct {
	Observable
	try       cancellableLocker
	observers []*ObserverFunc
	hasError  bool
	errValue  error
}

// Subject implements SubjectLike.
var _ SubjectLike = (*Subject)(nil)

// Next emits an value to the consumers of this Subject.
func (s *Subject) Next(val interface{}) {
	if s.try.Lock() {
		for _, ob := range s.observers {
			ob.Next(val)
		}
		s.try.Unlock()
	}
}

// Error emits an error Notification to the consumers of this Subject.
func (s *Subject) Error(err error) {
	if s.try.Lock() {
		observers := s.observers
		s.observers = nil
		s.hasError = true
		s.errValue = err

		s.try.CancelAndUnlock()

		for _, ob := range observers {
			ob.Error(err)
		}
	}
}

// Complete emits a complete Notification to the consumers of this Subject.
func (s *Subject) Complete() {
	if s.try.Lock() {
		observers := s.observers
		s.observers = nil

		s.try.CancelAndUnlock()

		for _, ob := range observers {
			ob.Complete()
		}
	}
}

// Subscribe adds a consumer to this Subject.
func (s *Subject) Subscribe(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
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

		s.try.Unlock()

		return ctx, cancel
	}

	if s.hasError {
		ob.Error(s.errValue)
	} else {
		ob.Complete()
	}

	return canceledCtx, noopFunc
}

// NewSubject returns a new Subject.
func NewSubject() *Subject {
	s := &Subject{}
	s.Op = OperatorFunc(s.Subscribe)
	return s
}
