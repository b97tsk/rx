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

// NewSubject creates a new Subject.
func NewSubject() Subject {
	s := new(subject)
	return Subject{
		Observable: Observable{}.Lift(s.call),
		Observer:   s.notify,
	}
}

type subject struct {
	try       cancellableLocker
	observers observerList
	err       error
}

func (s *subject) notify(t Notification) {
	if s.try.Lock() {
		switch {
		case t.HasValue:
			observers, releaseRef := s.observers.AddRef()

			s.try.Unlock()

			for _, sink := range observers {
				sink.Notify(t)
			}

			releaseRef()

		case t.HasError:
			observers := s.observers.Swap(nil)
			s.err = t.Value.(error)

			s.try.CancelAndUnlock()

			for _, sink := range observers {
				sink.Notify(t)
			}

		default:
			observers := s.observers.Swap(nil)

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

		observer := Mutex(Finally(sink, cancel))
		s.observers.Append(&observer)

		go func() {
			<-ctx.Done()
			if s.try.Lock() {
				s.observers.Remove(&observer)
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
