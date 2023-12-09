package rx

import (
	"context"
)

// ToSlice collects all the values emitted by the source Observable,
// and then emits them as a slice when the source completes.
func ToSlice[T any]() Operator[T, []T] {
	return NewOperator(toSlice[T])
}

func toSlice[T any](source Observable[T]) Observable[[]T] {
	return func(ctx context.Context, sink Observer[[]T]) {
		var s []T

		source.Subscribe(ctx, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				s = append(s, n.Value)
			case KindError:
				sink.Error(n.Error)
			case KindComplete:
				sink.Next(s)
				sink.Complete()
			}
		})
	}
}
