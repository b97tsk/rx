package rx

import (
	"context"
)

// Take emits only the first count values emitted by the source Observable.
func Take[T any](count int) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			if count <= 0 {
				return Empty[T]()
			}

			return takeObservable[T]{source, count}.Subscribe
		},
	)
}

type takeObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs takeObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var noop bool

	count := obs.Count

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if noop {
			return
		}

		sink(n)

		if n.HasValue {
			count--

			if count == 0 {
				noop = true

				sink.Complete()
			}
		}
	})
}
