package rx

import (
	"context"
)

// Last emits only the last value emitted by the source Observable.
// If the source turns out to be empty, it emits an error notification
// of ErrEmpty.
func Last[T any]() Operator[T, T] {
	return NewOperator(last[T])
}

func last[T any](source Observable[T]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		var (
			value     T
			haveValue bool
		)

		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case n.HasValue:
				value = n.Value
				haveValue = true
			case n.HasError:
				sink(n)
			default:
				if haveValue {
					sink.Next(value)
					sink.Complete()
				} else {
					sink.Error(ErrEmpty)
				}
			}
		})
	}
}

// LastOrElse emits only the last value emitted by the source Observable.
// If the source turns out to be empty, it emits a specified default value.
func LastOrElse[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				var (
					value     T
					haveValue bool
				)

				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						value = n.Value
						haveValue = true
					case n.HasError:
						sink(n)
					default:
						if haveValue {
							sink.Next(value)
						} else {
							sink.Next(def)
						}

						sink(n)
					}
				})
			}
		},
	)
}
