package rx

import (
	"context"
)

// Subject is a special type of Observable that allows values to be multicasted
// to many Observers.
type Subject struct {
	Observable
	Observer
}

type subject struct {
	try       cancellableLocker
	observers []*Observer
	err       error
}

func (s *subject) notify(t Notification) {
	if s.try.Lock() {
		switch {
		case t.HasValue:
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

func (s *subject) call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

		s.try.Unlock()
		return ctx, cancel
	}

	if s.err != nil {
		sink.Error(s.err)
	} else {
		sink.Complete()
	}

	return Done()
}

// NewSubject returns a new Subject.
func NewSubject() Subject {
	s := new(subject)
	return Subject{
		Observable: Observable{}.Lift(s.call),
		Observer:   s.notify,
	}
}
