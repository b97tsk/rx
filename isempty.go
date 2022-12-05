package rx

import (
	"context"
)

// IsEmpty emits a boolean to indicate whether the source emits no items.
func IsEmpty[T any]() Operator[T, bool] {
	return NewOperator(isEmpty[T])
}

func isEmpty[T any](source Observable[T]) Observable[bool] {
	return func(ctx context.Context, sink Observer[bool]) {
		ctx, cancel := context.WithCancel(ctx)

		var noop bool

		source.Subscribe(ctx, func(n Notification[T]) {
			if noop {
				return
			}

			noop = true

			cancel()

			switch {
			case n.HasValue:
				sink.Next(false)
				sink.Complete()
			case n.HasError:
				sink.Error(n.Error)
			default:
				sink.Next(true)
				sink.Complete()
			}
		})
	}
}
