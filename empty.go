package rx

// Empty returns an Observable that emits no values and immediately completes.
func Empty[T any]() Observable[T] {
	return func(_ Context, o Observer[T]) {
		o.Complete()
	}
}

// Never returns an Observable that never emits anything, except
// when a context cancellation is detected, emits an error notification
// of whatever that context reports.
func Never[T any]() Observable[T] {
	return func(c Context, o Observer[T]) {
		if c.Done() != nil {
			c.AfterFunc(func() {
				o.Error(c.Err())
			})
		}
	}
}

// Throw creates an Observable that emits no values and immediately emits
// an error notification of err.
func Throw[T any](err error) Observable[T] {
	return func(_ Context, o Observer[T]) {
		o.Error(err)
	}
}

// Oops creates an Observable that emits no values and immediately emits
// an error notification of ErrOops, after calling panic(v).
func Oops[T any](v any) Observable[T] {
	return func(_ Context, o Observer[T]) {
		defer o.Error(ErrOops)
		panic(v)
	}
}
