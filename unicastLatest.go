package rx

import (
	"context"
	"sync"
)

// UnicastLatest returns a Double whose Observable part only takes care of
// one single Observer (the latest one subscribes to it), which will receive
// emissions from Double's Observer part; the previous one will immediately
// receive an ErrDropped.
func UnicastLatest() Double {
	d := &unicastLatest{}

	return Double{
		Observable: d.subscribe,
		Observer:   d.sink,
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

func (d *unicastLatest) sink(t Notification) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch {
	case d.err != nil:
		break

	case t.HasValue:
		if ctx := d.obs.ctx; ctx != nil {
			if ctx.Err() != nil {
				d.obs.ctx, d.obs.sink = nil, nil
			} else {
				d.obs.sink(t)
			}
		}

	default:
		d.err = errCompleted

		if t.HasError {
			d.err = t.Error

			if d.err == nil {
				d.err = errNil
			}
		}

		obs := d.obs

		d.obs.ctx, d.obs.sink = nil, nil

		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			obs.sink(t)
		}
	}
}

func (d *unicastLatest) subscribe(ctx context.Context, sink Observer) {
	d.mu.Lock()

	err := d.err
	if err == nil {
		obs := d.obs

		d.obs.ctx, d.obs.sink = ctx, sink

		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			err = ErrDropped
			sink = obs.sink
		}
	}

	d.mu.Unlock()

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
