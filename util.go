package rx

import (
	"context"
)

type observables[T any] []Observable[T]

func identity[T any](v T) T { return v }

func subscribeToChan[T any](ctx context.Context, obs Observable[T], c chan<- Notification[T]) {
	done := ctx.Done()
	obs.Subscribe(ctx, func(n Notification[T]) {
		select {
		case <-done:
		case c <- n:
		}
	})
}
