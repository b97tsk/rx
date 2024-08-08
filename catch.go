package rx

// Catch mirrors the source or switches to another Observable, returned from
// a call to selector, if the source emits an [Error] notification.
func Catch[T any](selector func(err error) Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext, KindComplete, KindStop:
						o.Emit(n)
					case KindError:
						ob := Try11(selector, n.Error, func() { o.Stop(ErrOops) })
						ob.Subscribe(c, o)
					}
				})
			}
		},
	)
}

// OnErrorResumeWith mirrors the source or switches to specified Observable
// if the source emits an [Error] notification.
func OnErrorResumeWith[T any](ob Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext, KindComplete, KindStop:
						o.Emit(n)
					case KindError:
						ob.Subscribe(c, o)
					}
				})
			}
		},
	)
}

// OnErrorComplete mirrors the source Observable, or completes if the source
// emits an [Error] notification.
func OnErrorComplete[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext, KindComplete, KindStop:
						o.Emit(n)
					case KindError:
						o.Complete()
					}
				})
			}
		},
	)
}
