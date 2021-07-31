package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/internal/ctxutil"
)

// Multicast returns a Double whose Observable part takes care of all
// Observers that subscribes to it, which will receive emissions from
// Double's Observer part.
func Multicast() Double {
	d := &multicast{}

	return Double{
		Observable: d.subscribe,
		Observer:   d.sink,
	}
}

type multicast struct {
	mu  sync.Mutex
	err error
	lst observerList
	cws ctxutil.ContextWaitService
}

func (d *multicast) sink(t Notification) {
	d.mu.Lock()

	switch {
	case d.err != nil:
		d.mu.Unlock()

	case t.HasValue:
		lst := d.lst.Clone()
		defer lst.Release()

		d.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList

		d.lst.Swap(&lst)

		d.err = errCompleted

		if t.HasError {
			d.err = t.Error

			if d.err == nil {
				d.err = errNil
			}
		}

		d.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (d *multicast) subscribe(ctx context.Context, sink Observer) {
	d.mu.Lock()

	err := d.err
	if err == nil {
		ctx, cancel := context.WithCancel(ctx)

		observer := sink.WithCancel(cancel).MutexContext(ctx)

		d.lst.Append(&observer)

		d.cws.Submit(ctx, func() {
			d.mu.Lock()
			d.lst.Remove(&observer)
			d.mu.Unlock()
		})
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
