package rx

import (
	"context"
)

// Materialize represents all of the Notifications from the source Observable
// as values, and then completes.
func Materialize[T any]() Operator[T, Notification[T]] {
	return NewOperator(materialize[T])
}

func materialize[T any](source Observable[T]) Observable[Notification[T]] {
	return func(ctx context.Context, sink Observer[Notification[T]]) {
		source.Subscribe(ctx, func(n Notification[T]) {
			sink.Next(n)

			if !n.HasValue {
				sink.Complete()
			}
		})
	}
}
