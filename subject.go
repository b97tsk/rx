package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/x/misc"
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
	return Subject{Create(s.subscribe), s.sink}
}

type subject struct {
	mux sync.Mutex
	lst observerList
	cws misc.ContextWaitService
	err error
}

func (s *subject) sink(t Notification) {
	s.mux.Lock()
	switch {
	case s.err != nil:
		s.mux.Unlock()

	case t.HasValue:
		lst := s.lst.Clone()
		defer lst.Release()

		s.mux.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList
		s.lst.Swap(&lst)

		if t.HasError {
			s.err = t.Error
		} else {
			s.err = Completed
		}

		s.mux.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (s *subject) subscribe(ctx context.Context, sink Observer) {
	s.mux.Lock()

	err := s.err
	if err == nil {
		observer := Mutex(sink)
		s.lst.Append(&observer)

		finalize := func() {
			s.mux.Lock()
			s.lst.Remove(&observer)
			s.mux.Unlock()
		}

		for s.cws == nil || !s.cws.Submit(ctx, finalize) {
			s.cws = misc.NewContextWaitService()
		}
	}

	s.mux.Unlock()

	if err != nil {
		if err != Completed {
			sink.Error(err)
		} else {
			sink.Complete()
		}
	}
}
