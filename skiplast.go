package rx

import (
	"context"
)

// SkipLast skips the last count values emitted by the source Observable.
func SkipLast[T any](count int) Operator[T, T] {
	if count <= 0 {
		return AsOperator(identity[Observable[T]])
	}

	return skipLast[T](count)
}

func skipLast[T any](count int) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return skipLastObservable[T]{source, count}.Subscribe
		},
	)
}

type skipLastObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs skipLastObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	buffer := make([]T, obs.Count)
	bufferSize := obs.Count

	var index, count int

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			if count < bufferSize {
				count++
			} else {
				sink.Next(buffer[index])
			}

			buffer[index] = n.Value
			index = (index + 1) % bufferSize

		default:
			sink(n)
		}
	})
}
