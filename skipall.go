package rx

import (
	"context"
)

// SkipAll skips all values emitted by the source Observable.
func SkipAll[T any]() Operator[T, T] {
	return NewOperator(skipAll[T])
}

func skipAll[T any](source Observable[T]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		source.Subscribe(ctx, func(n Notification[T]) {
			if !n.HasValue {
				sink(n)
			}
		})
	}
}