package rx

// Do mirrors the source Observable, passing emissions to tap before
// each emission.
func Do[T any](tap Observer[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						tap(n)
						sink(n)
					case KindError, KindComplete:
						Try1(tap, n, func() { sink.Error(ErrOops) })
						sink(n)
					}
				})
			}
		},
	)
}

// DoOnNext mirrors the source Observable, passing values to f before
// each value emission.
func DoOnNext[T any](f func(v T)) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						f(n.Value)
					}
					sink(n)
				})
			}
		},
	)
}

// DoOnError mirrors the source Observable, and calls f when the source emits
// a notification of error.
func DoOnError[T any](f func(err error)) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindError {
						Try1(f, n.Error, func() { sink.Error(ErrOops) })
					}
					sink(n)
				})
			}
		},
	)
}

// DoOnComplete mirrors the source Observable, and calls f when the source
// completes.
func DoOnComplete[T any](f func()) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindComplete {
						Try0(f, func() { sink.Error(ErrOops) })
					}
					sink(n)
				})
			}
		},
	)
}

// DoOnTermination mirrors the source Observable, and calls f when the source
// emits a notification of error or completion.
func DoOnTermination[T any](f func()) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, sink.DoOnTermination(f))
			}
		},
	)
}
