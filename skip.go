package rx

import (
	"context"
)

// Skip skips the first count items emitted by the source Observable.
func Skip[T any](count int) Operator[T, T] {
	if count <= 0 {
		return AsOperator(identity[Observable[T]])
	}

	return skip[T](count)
}

func skip[T any](count int) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return skipObservable[T]{source, count}.Subscribe
		},
	)
}

type skipObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs skipObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	var taking bool

	count := obs.Count

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if taking || !n.HasValue {
			sink(n)
			return
		}

		count--
		taking = count == 0
	})
}
