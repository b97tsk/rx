package rx

import (
	"context"
	"sync"
)

// Complete is a special error that represents a complete subscription.
var Complete error = new(completeErr)

type completeErr int

func (*completeErr) Error() string {
	return "Complete"
}

// A Context is a context of a subscription.
type Context struct {
	context.Context
	mux    sync.Mutex
	err    error
	cancel context.CancelFunc
}

// NewContext creates a new Context based on an existing context.
func NewContext(ctx context.Context) *Context {
	ctx, cancel := context.WithCancel(ctx)
	return &Context{
		Context: ctx,
		cancel:  cancel,
	}
}

// AtLast is a helper method that is suitable to use with DoAtLast function.
func (c *Context) AtLast(t Notification) {
	switch {
	case t.HasValue:
		panic("rx.Context: AtLast() called with a NEXT notification")
	case t.HasError:
		c.Unsubscribe(t.Error)
	default:
		c.Unsubscribe(Complete)
	}
}

// Cancel cancels this context.
func (c *Context) Cancel() {
	c.Unsubscribe(context.Canceled)
}

// Unsubscribe cancels this context, and can also specify the returned
// value of Err() method.
func (c *Context) Unsubscribe(err error) {
	c.mux.Lock()
	if c.err == nil {
		c.err = err
		c.cancel()
	}
	c.mux.Unlock()
}

func (c *Context) Err() error {
	c.mux.Lock()
	err := c.err
	if err == nil {
		err = c.Context.Err()
		c.err = err
	}
	c.mux.Unlock()
	return err
}
