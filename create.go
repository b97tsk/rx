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
		sink = DoAtLast(sink, k.atLast)
		subscribe(k, sink)
		return k, cancel
	}
}

type kontext struct {
	context.Context
	mux    sync.Mutex
	err    error
	cancel context.CancelFunc
}

func (c *kontext) atLast(t Notification) {
	switch {
	case t.HasValue:
		panic("kontext: atLast() called with a NEXT notification")
	case t.HasError:
		c.unsubscribe(t.Error)
	default:
		c.unsubscribe(Complete)
	}
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
