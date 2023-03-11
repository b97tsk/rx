package rx

import (
	"context"
)

// First emits only the first value emitted by the source Observable.
// If the source turns out to be empty, it emits an error notification
// of ErrEmpty.
func First[T any]() Operator[T, T] {
	return NewOperator(first[T])
}

func first[T any](source Observable[T]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
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
				sink(n)
				sink.Complete()
			case n.HasError:
				sink(n)
			default:
				sink.Error(ErrEmpty)
			}
		})
	}
}

// FirstOrElse emits only the first value emitted by the source Observable.
// If the source turns out to be empty, it emits a specified default value.
func FirstOrElse[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
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
						sink(n)
						sink.Complete()
					case n.HasError:
						sink(n)
					default:
						sink.Next(def)
						sink(n)
					}
				})
			}
		},
	)
}
