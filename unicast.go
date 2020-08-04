package rx

import (
	"context"
	"sync"
)

// Unicast returns a Double whose Observable part only takes care of one single
// Observer (the first one subscribes to it), which will receive emissions from
// Double's Observer part; the others will immediately receive an ErrDropped.
func Unicast() Double {
	d := &unicast{}
	return Double{
		Observable: d.subscribe,
		Observer:   d.sink,
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

func (d *unicast) sink(t Notification) {
	d.mu.Lock()
	defer d.mu.Unlock()
	switch {
	case d.err != nil:
	case t.HasValue:
		if ctx := d.obs.ctx; ctx != nil {
			if ctx.Err() != nil {
				d.obs.ctx, d.obs.sink = nil, Noop
			} else {
				d.obs.sink(t)
			}
		}
	default:
		d.err = errCompleted
		if t.HasError {
			d.err = t.Error
		}
		obs := d.obs
		d.obs.ctx, d.obs.sink = nil, nil
		if ctx := obs.ctx; ctx != nil && ctx.Err() == nil {
			obs.sink(t)
		}
	}
}

func (d *unicast) subscribe(ctx context.Context, sink Observer) {
	d.mu.Lock()
	err := d.err
	if err == nil {
		if d.obs.sink == nil {
			d.obs.ctx, d.obs.sink = ctx, sink
		} else {
			err = ErrDropped
		}
	}
	d.mu.Unlock()
	if err != nil {
		if err != errCompleted {
			sink.Error(err)
		} else {
			sink.Complete()
		}
	}
}
