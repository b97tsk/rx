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
	s.lock = make(chan struct{}, 1)
	s.lock <- struct{}{}
	return Subject{
		Observable: Observable{}.Lift(s.call),
		Observer:   s.notify,
	}
}

type subject struct {
	lock      chan struct{}
	observers observerList
	err       error
}

func (s *subject) notify(t Notification) {
	if _, ok := <-s.lock; ok {
		switch {
		case t.HasValue:
			observers, releaseRef := s.observers.AddRef()

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

func (s *subject) call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

		s.lock <- struct{}{}
		return ctx, cancel
	}

	if s.err != nil {
		sink.Error(s.err)
	} else {
		sink.Complete()
	}

	return Done()
}
