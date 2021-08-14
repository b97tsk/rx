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
	u := &unicastLatest{}

	return Subject{
		Observable: u.subscribe,
		Observer:   u.sink,
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

func (u *unicastLatest) sink(t Notification) {
	u.mu.Lock()
	defer u.mu.Unlock()

	switch {
	case u.err != nil:
		break

	case t.HasValue:
		if ctx := u.obs.ctx; ctx != nil {
			if ctx.Err() != nil {
				u.obs.ctx, u.obs.sink = nil, nil
			} else {
				u.obs.sink(t)
			}
		}

	default:
		u.err = errCompleted

		if t.HasError {
			u.err = t.Error

			if u.err == nil {
				u.err = errNil
			}
		}

		obs := u.obs

		u.obs.ctx, u.obs.sink = nil, nil

		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			obs.sink(t)
		}
	}
}

func (u *unicastLatest) subscribe(ctx context.Context, sink Observer) {
	u.mu.Lock()

	err := u.err
	if err == nil {
		obs := u.obs

		u.obs.ctx, u.obs.sink = ctx, sink

		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			err = ErrDropped
			sink = obs.sink
		}
	}

	u.mu.Unlock()

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
