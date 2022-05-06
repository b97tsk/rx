package rx

import (
	"context"
)

// Map applies a given projection function to each value emitted by
// the source Observable, then emits the resulting values.
func Map[T, R any](proj func(v T) R) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return map1(proj)
}

func map1[T, R any](proj func(v T) R) Operator[T, R] {
	return AsOperator(
		func(source Observable[T]) Observable[R] {
			return func(ctx context.Context, sink Observer[R]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						sink.Next(proj(n.Value))
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Complete()
					}
				})
			}
		},
	)
}

// MapTo emits the given constant value on the output Observable every time
// the source Observable emits a value.
func MapTo[T, R any](v R) Operator[T, R] {
	return AsOperator(
		func(source Observable[T]) Observable[R] {
			return func(ctx context.Context, sink Observer[R]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						sink.Next(v)
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Complete()
					}
				})
			}
		},
	)
}
