package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/internal/ctxutil"
)

// Multicast returns a Subject whose Observable part takes care of all
// Observers that subscribes to it, which will receive emissions from
// Subject's Observer part.
func Multicast() Subject {
	s := &multicast{}

	return Subject{
		Observable: s.subscribe,
		Observer:   s.sink,
	}
}

type multicast struct {
	mu  sync.Mutex
	err error
	lst observerList
	cws ctxutil.ContextWaitService
}

func (s *multicast) sink(t Notification) {
	s.mu.Lock()

	switch {
	case s.err != nil:
		s.mu.Unlock()

	case t.HasValue:
		lst := s.lst.Clone()
		defer lst.Release()

		s.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList

		s.lst.Swap(&lst)

		s.err = errCompleted

		if t.HasError {
			s.err = t.Error

			if s.err == nil {
				s.err = errNil
			}
		}

		s.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (s *multicast) subscribe(ctx context.Context, sink Observer) {
	s.mu.Lock()

	err := s.err
	if err == nil {
		ctx, cancel := context.WithCancel(ctx)

		observer := sink.WithCancel(cancel).MutexContext(ctx)

		s.lst.Append(&observer)

		s.cws.Submit(ctx, func() {
			s.mu.Lock()
			s.lst.Remove(&observer)
			s.mu.Unlock()
		})
	}

	s.mu.Unlock()

	if err != nil {
		if err == errCompleted {
			sink.Complete()
			return
		}

		if err == errNil {
			err = nil
		}

		sink.Error(err)
	}
}
