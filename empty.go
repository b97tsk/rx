package rx

// Empty creates an [Observable] that emits no values and immediately
// completes.
func Empty[T any]() Observable[T] {
	return func(_ Context, o Observer[T]) {
		o.Complete()
	}
}

// Never creates an [Observable] that never emits anything, except when
// a [Context] cancellation is detected, emits a [Stop] notification of
// whatever [Context.Cause] returns.
func Never[T any]() Observable[T] {
	return func(c Context, o Observer[T]) {
		if c.Done() != nil {
			c.AfterFunc(func() {
				o.Stop(c.Cause())
			})
		}
	}
}

// Throw creates an [Observable] that emits no values and immediately emits
// an [Error] notification of err.
func Throw[T any](err error) Observable[T] {
	return func(_ Context, o Observer[T]) {
		o.Error(err)
	}
}

// Oops creates an [Observable] that emits no values and immediately emits
// a [Stop] notification of [ErrOops], after running panic(v).
func Oops[T any](v any) Observable[T] {
	return func(_ Context, o Observer[T]) {
		defer o.Stop(ErrOops)
		panic(v)
	}
}
