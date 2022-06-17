package rx

import (
	"context"
)

// ToSlice collects all the values emitted by the source Observable, and then
// emits them as a slice when the source completes.
func ToSlice[T any]() Operator[T, []T] {
	return AsOperator(toSlice[T])
}

func toSlice[T any](source Observable[T]) Observable[[]T] {
	return func(ctx context.Context, sink Observer[[]T]) {
		var s []T

		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case n.HasValue:
				s = append(s, n.Value)
			case n.HasError:
				sink.Error(n.Error)
			default:
				sink.Next(s)
				sink.Complete()
			}
		})
	}
}
