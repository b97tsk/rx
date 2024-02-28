package rx

import (
	"context"
	"sync"
	"time"
)

// A CancelFunc tells an operation to abandon its work.
// A CancelFunc does not wait for the work to stop.
// A CancelFunc may be called by multiple goroutines simultaneously.
// After the first call, subsequent calls to a CancelFunc do nothing.
type CancelFunc = context.CancelFunc

// A Context carries a [context.Context], an optional [sync.WaitGroup], and
// an optional panic handler.
type Context struct {
	Context      context.Context
	WaitGroup    *sync.WaitGroup
	PanicHandler func(v any)
}

// NewBackgroundContext returns NewContext(context.Background()).
func NewBackgroundContext() Context {
	return NewContext(context.Background())
}

// NewContext returns a [Context] with Context field set to ctx.
func NewContext(ctx context.Context) Context {
	return Context{Context: ctx}
}

// Done returns c.Context.Done().
func (c Context) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err returns c.Context.Err().
func (c Context) Err() error {
	return c.Context.Err()
}

// AfterFunc arranges to call f in its own goroutine after c is done
// (cancelled or timed out).
// If c is already done, AfterFunc calls f immediately in its own goroutine.
//
// Calling the returned stop function stops the association of c with f.
// It returns true if the call stopped f from being run.
// If stop returns false, either the context is done and f has been started
// in its own goroutine; or f was already stopped.
// The stop function does not wait for f to complete before returning.
//
// Internally, f is wrapped with [Context.PreAsyncCall] before being passed
// to [context.AfterFunc].
func (c Context) AfterFunc(f func()) (stop func() bool) {
	return context.AfterFunc(c.Context, c.PreAsyncCall(f))
}

// WithCancel returns a copy of c with a new Done channel.
// The returned context's Done channel is closed when the returned cancel
// function is called or when c's Done channel is closed, whichever happens
// first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this context complete.
func (c Context) WithCancel() (Context, CancelFunc) {
	ctx, cancel := context.WithCancel(c.Context)
	c.Context = ctx
	return c, cancel
}

// WithDeadline returns a copy of c with the deadline adjusted to be no later
// than d. If c's deadline is already earlier than d, c.WithDeadline(d) is
// semantically equivalent to c. The returned context's Done channel is closed
// when the deadline expires, when the returned cancel function is called, or
// when c's Done channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this context complete.
func (c Context) WithDeadline(d time.Time) (Context, CancelFunc) {
	ctx, cancel := context.WithDeadline(c.Context, d)
	c.Context = ctx
	return c, cancel
}

// WithTimeout returns c.WithDeadline(time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this context complete.
func (c Context) WithTimeout(timeout time.Duration) (Context, CancelFunc) {
	return c.WithDeadline(time.Now().Add(timeout))
}

// WithWaitGroup returns a copy of c with WaitGroup field set to wg.
func (c Context) WithWaitGroup(wg *sync.WaitGroup) Context {
	c.WaitGroup = wg
	return c
}

// WithPanicHandler returns a copy of c with PanicHandler field set to f.
func (c Context) WithPanicHandler(f func(v any)) Context {
	c.PanicHandler = f
	return c
}

// Go calls f in a goroutine.
//
// Internally, f is wrapped with [Context.PreAsyncCall] before being run by
// the built-in go statement.
func (c Context) Go(f func()) {
	go c.PreAsyncCall(f)()
}

// PreAsyncCall increases c.WaitGroup's counter, if c.WaitGroup is not nil,
// and returns a function that calls f.
//
// If c.WaitGroup is not nil, the function returned decreases c.WaitGroup's
// counter when f returns.
//
// If f panics and c.PanicHandler is not nil, the function returned calls
// c.PanicHandler with a value returned by the built-in recover function.
//
// PreAsyncCall is usually called before starting an asynchronous operation,
// the caller then calls the function returned in that asynchronous operation.
// The function passed to PreAsyncCall is what the caller would do in that
// asynchronous operation.
func (c Context) PreAsyncCall(f func()) func() {
	if c.WaitGroup != nil {
		c.WaitGroup.Add(1)
	}
	return func() {
		defer func() {
			if c.WaitGroup != nil {
				c.WaitGroup.Done()
			}
			if c.PanicHandler != nil {
				if v := recover(); v != nil {
					c.PanicHandler(v)
				}
			}
		}()
		f()
	}
}
