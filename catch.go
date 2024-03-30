package rx

// Catch mirrors the source or switches to another Observable, returned from
// a call to selector, if the source emits a notification of error.
//
// Catch does not catch context cancellations.
func Catch[T any](selector func(err error) Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						sink(n)

					case KindError:
						select {
						default:
						case <-c.Done():
							sink.Error(c.Err())
							return
						}

						obs := Try11(selector, n.Error, func() { sink.Error(ErrOops) })
						obs.Subscribe(c, sink)

					case KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}

// OnErrorResumeWith mirrors the source or specified Observable if the source
// emits a notification of error.
//
// OnErrorResumeWith does not resume after context cancellation.
func OnErrorResumeWith[T any](obs Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						sink(n)

					case KindError:
						select {
						default:
						case <-c.Done():
							sink.Error(c.Err())
							return
						}

						obs.Subscribe(c, sink)

					case KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}

// OnErrorComplete mirrors the source Observable, or completes if the source
// emits a notification of error.
//
// OnErrorComplete does not complete after context cancellation.
func OnErrorComplete[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						sink(n)

					case KindError:
						select {
						default:
						case <-c.Done():
							sink.Error(c.Err())
							return
						}

						sink.Complete()

					case KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}
