package rx

import (
	"context"
)

// Subject is a special type of Observable that allows values to be multicasted
// to many Observers.
type Subject struct {
	Observer
	Observable
	try       cancellableLocker
	observers []*Observer
	errValue  error
}

// Next emits an value to the consumers of this Subject.
func (s *Subject) Next(val interface{}) {
	if s.try.Lock() {
		defer s.try.Unlock()
		for _, ob := range s.observers {
			ob.Next(val)
		}
	}
}

// Error emits an error Notification to the consumers of this Subject.
func (s *Subject) Error(err error) {
	if s.try.Lock() {
		observers := s.observers
		s.observers = nil
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

		return ctx, cancel
	}

	if s.errValue != nil {
		ob.Error(s.errValue)
	} else {
		ob.Complete()
	}

	return canceledCtx, noopFunc
}

// NewSubject returns a new Subject.
func NewSubject() *Subject {
	s := new(Subject)
	s.Observer = func(t Notification) {
		switch {
		case t.HasValue:
			s.Next(t.Value)
		case t.HasError:
			s.Error(t.Value.(error))
		default:
			s.Complete()
		}
	}
	s.Observable = s.Observable.Lift(
		func(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
			return s.Subscribe(ctx, ob)
		},
	)
	return s
}
