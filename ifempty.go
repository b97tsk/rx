package rx

import (
	"context"
)

// JustIfEmpty mirrors the source Observable, or emits given values
// if the source completes without emitting any value.
func JustIfEmpty[T any](s ...T) Operator[T, T] {
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
							done := ctx.Done()

							for _, v := range s {
								select {
								default:
								case <-done:
									sink.Error(ctx.Err())
									return
								}

								sink.Next(v)
							}
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
