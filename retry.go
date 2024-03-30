package rx

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

func (obs retryObservable[T]) Subscribe(c Context, sink Observer[T]) {
	var observer Observer[T]

	done := c.Done()

	subscribeToSource := resistReentrance(func() {
		select {
		default:
		case <-done:
			sink.Error(c.Err())
			return
		}

		obs.Source.Subscribe(c, observer)
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
