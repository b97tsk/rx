package rx

import (
	"context"
)

// DefaultIfEmpty mirrors the source Observable, emits a given value if
// the source completes without emitting any value.
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

// ThrowIfEmpty mirrors the source Observable, throws ErrEmpty if the source
// completes without emitting any value.
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
	if obs == nil {
		panic("obs == nil")
	}

	return switchIfEmpty(obs)
}

func switchIfEmpty[T any](obs Observable[T]) Operator[T, T] {
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
