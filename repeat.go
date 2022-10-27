package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/norec"
)

// RepeatForever repeats the stream of items emitted by the source Observable
// forever.
//
// RepeatForever does not repeat on context cancellation.
//
func RepeatForever[T any]() Operator[T, T] {
	return Repeat[T](-1)
}

// Repeat repeats the stream of items emitted by the source Observable at
// most count times.
//
// Repeat(0) results in an empty Observable; Repeat(1) is a no-op.
//
// Repeat does not repeat on context cancellation.
//
func Repeat[T any](count int) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			if count == 0 {
				return Empty[T]()
			}

			if count == 1 {
				return source
			}

			if count > 0 {
				count--
			}

			return repeatObservable[T]{source, count}.Subscribe
		},
	)
}

type repeatObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs repeatObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	var observer Observer[T]

	subscribeToSource := norec.Wrap(func() {
		if err := ctx.Err(); err != nil {
			sink.Error(err)
			return
		}

		obs.Source.Subscribe(ctx, observer)
	})

	count := obs.Count

	observer = func(n Notification[T]) {
		if n.HasValue || n.HasError || count == 0 {
			sink(n)
			return
		}

		if count > 0 {
			count--
		}

		subscribeToSource()
	}

	subscribeToSource()
}
