package rx

// Do mirrors the source Observable, passing emissions to tap before
// each emission.
func Do[T any](tap Observer[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						tap(n)
						o.Emit(n)
					case KindError, KindComplete:
						Try1(tap, n, func() { o.Error(ErrOops) })
						o.Emit(n)
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
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						f(n.Value)
					}
					o.Emit(n)
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
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindError {
						Try1(f, n.Error, func() { o.Error(ErrOops) })
					}
					o.Emit(n)
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
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindComplete {
						Try0(f, func() { o.Error(ErrOops) })
					}
					o.Emit(n)
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
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, o.DoOnTermination(f))
			}
		},
	)
}
