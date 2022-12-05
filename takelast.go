package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// TakeLast emits only the last count values emitted by the source Observable.
func TakeLast[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count <= 0 {
				return Empty[T]()
			}

			return takeLastObservable[T]{source, count}.Subscribe
		},
	)
}

type takeLastObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs takeLastObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	var q queue.Queue[T]

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			if q.Len() == obs.Count {
				q.Pop()
			}

			q.Push(n.Value)

		case n.HasError:
			sink(n)

		default:
			for i, j := 0, q.Len(); i < j; i++ {
				if err := ctx.Err(); err != nil {
					sink.Error(err)
					return
				}

				sink.Next(q.At(i))
			}

			sink(n)
		}
	})
}
