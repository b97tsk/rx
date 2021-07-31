package rx

import (
	"context"
	"sync"
)

// UnicastLatest returns a Subject whose Observable part only takes care of
// one single Observer (the latest one subscribes to it), which will receive
// emissions from Subject's Observer part; the previous one will immediately
// receive an ErrDropped.
func UnicastLatest() Subject {
	s := &unicastLatest{}

	return Subject{
		Observable: s.subscribe,
		Observer:   s.sink,
	}
}

type unicastLatest struct {
	mu  sync.Mutex
	err error
	obs struct {
		ctx  context.Context
		sink Observer
	}
}

func (s *unicastLatest) sink(t Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch {
	case s.err != nil:
		break

	case t.HasValue:
		if ctx := s.obs.ctx; ctx != nil {
			if ctx.Err() != nil {
				s.obs.ctx, s.obs.sink = nil, nil
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

func (s *unicastLatest) subscribe(ctx context.Context, sink Observer) {
	s.mu.Lock()

	err := s.err
	if err == nil {
		obs := s.obs

		s.obs.ctx, s.obs.sink = ctx, sink

		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			err = ErrDropped
			sink = obs.sink
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
