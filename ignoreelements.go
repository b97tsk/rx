package rx

import (
	"context"
)

// IgnoreElements ignores all values emitted by the source Observable
// and only mirrors notifications of error or completion.
func IgnoreElements[T, R any]() Operator[T, R] {
	return NewOperator(ignoreElements[T, R])
}

func ignoreElements[T, R any](source Observable[T]) Observable[R] {
	return func(ctx context.Context, sink Observer[R]) {
		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case n.HasValue:
			case n.HasError:
				sink.Error(n.Error)
			default:
				sink.Complete()
			}
		})
	}
}
