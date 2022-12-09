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

func chanObserver[T any](c chan<- Notification[T], noop <-chan struct{}) Observer[T] {
	return func(n Notification[T]) {
		select {
		case c <- n:
		case <-noop:
		}
	}
}
