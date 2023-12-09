package rx

import (
	"context"
)

// IsEmpty emits a boolean to indicate whether the source emits no values.
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

			switch n.Kind {
			case KindNext, KindError, KindComplete:
				noop = true

				cancel()

				switch n.Kind {
				case KindNext:
					sink.Next(false)
					sink.Complete()
				case KindError:
					sink.Error(n.Error)
				case KindComplete:
					sink.Next(true)
					sink.Complete()
				}
			}
		})
	}
}
