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

// OnNext mirrors the source Observable, passing values to f before
// each value emission.
func OnNext[T any](f func(v T)) Operator[T, T] {
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

// OnComplete mirrors the source Observable, and calls f when the source
// completes.
func OnComplete[T any](f func()) Operator[T, T] {
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

// OnError mirrors the source Observable, and calls f when the source emits
// an error notification.
func OnError[T any](f func(err error)) Operator[T, T] {
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

// OnLastNotification mirrors the source Observable, and calls f when
// the source completes or emits an error notification.
func OnLastNotification[T any](f func()) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, sink.OnLastNotification(f))
			}
		},
	)
}
