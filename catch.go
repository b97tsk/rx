package rx

import (
	"context"
)

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
//
// Catch does not catch context cancellations.
func Catch[T any](selector func(err error) Observable[T]) Operator[T, T] {
	if selector == nil {
		panic("selector == nil")
	}

	return catch(selector)
}

func catch[T any](selector func(err error) Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						sink(n)

					case n.HasError:
						if err := ctx.Err(); err != nil {
							sink.Error(err)
							return
						}

						obs := selector(n.Error)
						obs.Subscribe(ctx, sink)

					default:
						sink(n)
					}
				})
			}
		},
	)
}

// OnErrorResumeWith mirrors the source or specified Observable if the source
// throws an error.
//
// OnErrorResumeWith does not resume on context cancellation.
func OnErrorResumeWith[T any](obs Observable[T]) Operator[T, T] {
	if obs == nil {
		panic("obs == nil")
	}

	return onErrorResumeWith(obs)
}

func onErrorResumeWith[T any](obs Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						sink(n)

					case n.HasError:
						if err := ctx.Err(); err != nil {
							sink.Error(err)
							return
						}

						obs.Subscribe(ctx, sink)

					default:
						sink(n)
					}
				})
			}
		},
	)
}

// OnErrorComplete mirrors the source Observable or completes if the source
// throws an error.
//
// OnErrorComplete does not complete on context cancellation.
func OnErrorComplete[T any]() Operator[T, T] {
	return NewOperator(onErrorComplete[T])
}

func onErrorComplete[T any](source Observable[T]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case n.HasValue:
				sink(n)

			case n.HasError:
				if err := ctx.Err(); err != nil {
					sink.Error(err)
					return
				}

				sink.Complete()

			default:
				sink(n)
			}
		})
	}
}
