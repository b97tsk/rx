package rx

// First emits only the first value emitted by the source [Observable].
// If the source turns out to be empty, First emits an [Error] notification
// of [ErrEmpty].
func First[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						o.Emit(n)
						noop = true
						o.Complete()
					case KindComplete:
						o.Error(ErrEmpty)
					case KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}

// FirstOrElse emits only the first value emitted by the source [Observable].
// If the source turns out to be empty, FirstOrElse emits a specified default
// value.
func FirstOrElse[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						o.Emit(n)
						noop = true
						o.Complete()
					case KindComplete:
						Try1(o, Next(def), func() { o.Stop(ErrOops) })
						o.Emit(n)
					case KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}
