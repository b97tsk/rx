package rx

// DefaultIfEmpty mirrors the source Observable, or emits given values
// if the source completes without emitting any value.
func DefaultIfEmpty[T any](s ...T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				haveValue := false

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						haveValue = true
					case KindError:
					case KindComplete:
						if !haveValue {
							done := c.Done()

							for _, v := range s {
								select {
								default:
								case <-done:
									o.Error(c.Cause())
									return
								}

								Try1(o, Next(v), func() { o.Error(ErrOops) })
							}
						}
					}

					o.Emit(n)
				})
			}
		},
	)
}

// ThrowIfEmpty mirrors the source Observable, or emits a notification of
// ErrEmpty if the source completes without emitting any value.
func ThrowIfEmpty[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				haveValue := false

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						haveValue = true
					case KindError:
					case KindComplete:
						if !haveValue {
							o.Error(ErrEmpty)
							return
						}
					}

					o.Emit(n)
				})
			}
		},
	)
}

// SwitchIfEmpty mirrors the source or specified Observable if the source
// completes without emitting any value.
func SwitchIfEmpty[T any](ob Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				haveValue := false

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						haveValue = true
					case KindError:
					case KindComplete:
						if !haveValue {
							ob.Subscribe(c, o)
							return
						}
					}

					o.Emit(n)
				})
			}
		},
	)
}
