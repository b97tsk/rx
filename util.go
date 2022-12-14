package rx

import (
	"context"
)

type observables[T any] []Observable[T]

type subscription struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func identity[T any](v T) T { return v }

// getErr returns ctx.Err().
//
// ctx.Err() is a bit slower than ctx.Done(), getErr avoids calling it
// when possible.
func getErr(ctx context.Context) error {
	return getErrWithDoneChan(ctx, ctx.Done())
}

// getErrWithDoneChan returns ctx.Err() if done is closed; otherwise
// it returns nil.
func getErrWithDoneChan(ctx context.Context, done <-chan struct{}) error {
	select {
	case <-done:
		return ctx.Err()
	default:
		return nil
	}
}

func chanObserver[T any](c chan<- Notification[T], noop <-chan struct{}) Observer[T] {
	return func(n Notification[T]) {
		select {
		case c <- n:
		case <-noop:
		}
	}
}
