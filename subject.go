package rx

import (
	"context"
	"sync"
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
		Observable: Create(s.subscribe),
		Observer:   s.notify,
	}
}

type subject struct {
	mux       sync.Mutex
	observers observerList
	cws       contextWaitService
	err       error
}

func (s *subject) notify(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		observers, releaseRef := s.observers.AddRef()

		s.mux.Unlock()

		for _, sink := range observers {
			sink.Notify(t)
		}

		releaseRef()

	default:
		observers := s.observers.Swap(nil)
		if t.HasError {
			s.err = t.Error
		} else {
			s.err = Complete
		}
		s.mux.Unlock()

		for _, sink := range observers {
			sink.Notify(t)
		}
	}
}

func (s *subject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	if err := s.err; err != nil {
		if err != Complete {
			sink.Error(err)
		} else {
			sink.Complete()
		}
	} else {
		observer := Mutex(sink)
		s.observers.Append(&observer)

		finalize := func() {
			s.mux.Lock()
			s.observers.Remove(&observer)
			s.mux.Unlock()
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = newContextWaitService()
		}
	}

	s.mux.Unlock()
}
