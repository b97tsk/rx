package rx

import (
	"context"
)

// Flat flattens a higher-order Observable into a first-order Observable,
// by applying a flat function to the inner Observables.
func Flat[_ Observable[T], T any](f func(some ...Observable[T]) Observable[T]) Operator[Observable[T], T] {
	if f == nil {
		panic("f == nil")
	}

	return flat(f)
}

func flat[_ Observable[T], T any](f func(some ...Observable[T]) Observable[T]) Operator[Observable[T], T] {
	return NewOperator(
		func(source Observable[Observable[T]]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				var s []Observable[T]

				source.Subscribe(ctx, func(n Notification[Observable[T]]) {
					switch {
					case n.HasValue:
						s = append(s, n.Value)
					case n.HasError:
						sink.Error(n.Error)
					default:
						f(s...).Subscribe(ctx, sink)
					}
				})
			}
		},
	)
}
