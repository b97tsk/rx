package rx

import (
	"context"
)

// IgnoreElements ignores all values emitted by the source Observable.
//
// It's like [SkipAll], but it can also change the output Observable to be
// of another type.
func IgnoreElements[T, R any]() Operator[T, R] {
	return NewOperator(ignoreElements[T, R])
}

func ignoreElements[T, R any](source Observable[T]) Observable[R] {
	return func(ctx context.Context, sink Observer[R]) {
		source.Subscribe(ctx, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
			case KindError:
				sink.Error(n.Error)
			case KindComplete:
				sink.Complete()
			}
		})
	}
}
