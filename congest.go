package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Congest mirrors the source Observable, buffers emissions if the source
// emits too fast, and blocks the source if the buffer is full.
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

	sink = sink.WithCancel(cancel)

	cout := make(chan Notification[T])

	go func() {
		done := ctx.Done()

		for {
			select {
			case <-done:
				sink.Error(ctx.Err())
				return
			case n := <-cout:
				sink(n)

				if !n.HasValue {
					return
				}
			}
		}
	}()

	bufferSize := obs.BufferSize

	cin := make(chan Notification[T])

	go func() {
		var q queue.Queue[Notification[T]]

		done := ctx.Done()

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
			case <-done:
				return
			case n := <-in:
				q.Push(n)
			case out <- outv:
				q.Pop()
			}
		}
	}()

	subscribeToChan(ctx, obs.Source, cin)
}
