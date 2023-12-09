package rx

import (
	"context"
)

// RetryForever mirrors the source Observable, and resubscribes to the source
// whenever the source emits a notification of error.
//
// RetryForever does not retry after context cancellation.
func RetryForever[T any]() Operator[T, T] {
	return Retry[T](-1)
}

// Retry mirrors the source Observable, and resubscribes to the source
// when the source emits a notification of error, for a maximum of count
// resubscriptions.
//
// Retry(0) is a no-op.
//
// Retry does not retry after context cancellation.
func Retry[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count == 0 {
				return source
			}

			return retryObservable[T]{source, count}.Subscribe
		},
	)
}

type retryObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs retryObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	var observer Observer[T]

	done := ctx.Done()

	subscribeToSource := resistReentry(func() {
		select {
		default:
		case <-done:
			sink.Error(ctx.Err())
			return
		}

		obs.Source.Subscribe(ctx, observer)
	})

	count := obs.Count

	observer = func(n Notification[T]) {
		if n.Kind != KindError || count == 0 {
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
