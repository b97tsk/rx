package rx

// Empty returns an Observable that emits no values and immediately completes.
func Empty[T any]() Observable[T] {
	return empty[T]
}

func empty[T any](_ Context, sink Observer[T]) {
	sink.Complete()
}

// Never returns an Observable that never emits anything, except
// when a context cancellation is detected, emits an error notification
// of whatever that context reports.
func Never[T any]() Observable[T] {
	return never[T]
}

func never[T any](c Context, sink Observer[T]) {
	if c.Done() != nil {
		c.AfterFunc(func() {
			sink.Error(c.Err())
		})
	}
}

// Throw creates an Observable that emits no values and immediately emits
// an error notification of err.
func Throw[T any](err error) Observable[T] {
	return func(_ Context, sink Observer[T]) {
		sink.Error(err)
	}
}
