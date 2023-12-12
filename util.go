package rx

import (
	"context"
	"sync/atomic"
)

// sentinel provides a distinct value for atomic operations.
var sentinel = func() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx
}()

type observables[T any] []Observable[T]

func identity[T any](v T) T { return v }

// resistReentrance returns a function that calls f in a non-recursive way
// when the function returned is called recursively.
func resistReentrance(f func()) func() {
	var n atomic.Uint32

	return func() {
		if n.Add(1) > 1 {
			return
		}

		for {
			f()

			if n.Add(^uint32(0)) == 0 {
				break
			}
		}
	}
}

func channelObserver[T any](c chan<- Notification[T], noop <-chan struct{}) Observer[T] {
	return func(n Notification[T]) {
		select {
		case c <- n:
		case <-noop:
		}
	}
}
