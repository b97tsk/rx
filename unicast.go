package rx

import (
	"context"
	"sync"
)

// Unicast returns a Subject whose Observable part only takes care of one single
// Observer (the first one subscribes to it), which will receive emissions from
// Subject's Observer part; the others will immediately receive an ErrDropped.
func Unicast() Subject {
	s := &unicast{}

	return Subject{
		Observable: s.subscribe,
		Observer:   s.sink,
	}
}

type unicast struct {
	mu  sync.Mutex
	err error
	obs struct {
		ctx  context.Context
		sink Observer
	}
}

func (s *unicast) sink(t Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch {
	case s.err != nil:
		break

	case t.HasValue:
		if ctx := s.obs.ctx; ctx != nil {
			if ctx.Err() != nil {
				s.obs.ctx, s.obs.sink = nil, Noop
			} else {
				s.obs.sink(t)
			}
		}

	default:
		s.err = errCompleted

		if t.HasError {
			s.err = t.Error

			if s.err == nil {
				s.err = errNil
			}
		}

		obs := s.obs

		s.obs.ctx, s.obs.sink = nil, nil

		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			obs.sink(t)
		}
	}
}

func (s *unicast) subscribe(ctx context.Context, sink Observer) {
	s.mu.Lock()

	err := s.err
	if err == nil {
		if s.obs.sink == nil {
			s.obs.ctx, s.obs.sink = ctx, sink
		} else {
			err = ErrDropped
		}
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
