package rx

import (
	"context"
)

// Subject is a special type of Observable that allows values to be multicasted
// to many Observers.
type Subject interface {
	Observer
	AsObservable() Observable
	Subscribe(context.Context, Observer) (context.Context, context.CancelFunc)
}

type subjectImplement struct {
	try       cancellableLocker
	observers []*ObserverFunc
	hasError  bool
	errValue  error
}

func (s *subjectImplement) Next(val interface{}) {
	if s.try.Lock() {
		for _, ob := range s.observers {
			ob.Next(val)
		}
		s.try.Unlock()
	}
}

func (s *subjectImplement) Error(err error) {
	if s.try.Lock() {
		observers := s.observers
		s.observers = nil
		s.hasError = true
		s.errValue = err

		s.try.Cancel()

		for _, ob := range observers {
			ob.Error(err)
		}
	}
}

func (s *subjectImplement) Complete() {
	if s.try.Lock() {
		observers := s.observers
		s.observers = nil

		s.try.Cancel()

		for _, ob := range observers {
			ob.Complete()
		}
	}
}

func (s *subjectImplement) AsObservable() Observable {
	return Observable{OperatorFunc(s.Subscribe)}
}

func (s *subjectImplement) Subscribe(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
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
func NewSubject() Subject {
	return &subjectImplement{}
}
