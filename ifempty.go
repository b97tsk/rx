package rx

import (
	"context"
)

// DefaultIfEmpty mirrors the source Observable, or emits a given value
// if the source completes without emitting any value.
func DefaultIfEmpty[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				haveValue := false

				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						haveValue = true
					case n.HasError:
					default:
						if !haveValue {
							sink.Next(def)
						}
					}

					sink(n)
				})
			}
		},
	)
}

// ThrowIfEmpty mirrors the source Observable, or emits an error notification
// of ErrEmpty if the source completes without emitting any value.
func ThrowIfEmpty[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				haveValue := false

				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						haveValue = true
					case n.HasError:
					default:
						if !haveValue {
							sink.Error(ErrEmpty)
							return
						}
					}

					sink(n)
				})
			}
		},
	)
}

// SwitchIfEmpty mirrors the source or specified Observable if the source
// completes without emitting any value.
func SwitchIfEmpty[T any](obs Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				haveValue := false

				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						haveValue = true
					case n.HasError:
					default:
						if !haveValue {
							obs.Subscribe(ctx, sink)
							return
						}
					}

					sink(n)
				})
			}
		},
	)
}
