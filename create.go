package rx

import (
	"context"
	"sync"
)

// Complete is a special error that represents a complete subscription.
var Complete error = new(complete)

type complete int

func (*complete) Error() string {
	return "complete"
}

// Create creates a new Observable, that will execute the specified function
// when an Observer subscribes to it.
//
// It's the caller's responsibility to follow the rule of Observable that
// no more emissions pass to the sink Observer after an ERROR or COMPLETE
// emission passes to it.
func Create(subscribe func(context.Context, Observer)) Observable {
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancel(ctx)
		k := &kontext{Context: ctx, cancel: cancel}
		subscribe(k, func(t Notification) {
			sink(t)
			switch {
			case t.HasValue:
			case t.HasError:
				k.unsubscribe(t.Error)
			default:
				k.unsubscribe(Complete)
			}
		})
		return k, cancel
	}
}

type kontext struct {
	context.Context
	mux    sync.Mutex
	err    error
	cancel context.CancelFunc
}

func (c *kontext) unsubscribe(err error) {
	c.mux.Lock()
	c.err = err
	c.cancel()
	c.mux.Unlock()
}

func (c *kontext) Err() error {
	c.mux.Lock()
	err := c.err
	if err == nil {
		err = c.Context.Err()
		c.err = err
	}
	c.mux.Unlock()
	return err
}
