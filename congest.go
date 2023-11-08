package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Congest mirrors the source Observable, buffering emissions if the source
// emits too fast, and blocking the source if the buffer is full.
//
// Congest has no effect if bufferSize < 1.
func Congest[T any](bufferSize int) Operator[T, T] {
	if bufferSize < 1 {
		return NewOperator(identity[Observable[T]])
	}

	return congest[T](bufferSize)
}

func congest[T any](bufferSize int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return congestObservable[T]{source, bufferSize}.Subscribe
		},
	)
}

type congestObservable[T any] struct {
	Source     Observable[T]
	BufferSize int
}

func (obs congestObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	cout := make(chan Notification[T])

	wg := WaitGroupFromContext(ctx)

	wg.Go(func() {
		for n := range cout {
			sink(n)
		}
	})

	bufferSize := obs.BufferSize

	cin := make(chan Notification[T])
	noop := make(chan struct{})

	wg.Go(func() {
		var q queue.Queue[Notification[T]]

		for {
			var (
				in   <-chan Notification[T]
				out  chan<- Notification[T]
				outv Notification[T]
			)

			length := q.Len()

			if length < bufferSize {
				in = cin
			}

			if length > 0 {
				out, outv = cout, q.Front()
			}

			select {
			case n := <-in:
				q.Push(n)
			case out <- outv:
				q.Pop()

				if !outv.HasValue {
					close(out)
					close(noop)

					return
				}
			}
		}
	})

	obs.Source.Subscribe(ctx, chanObserver(cin, noop))
}
