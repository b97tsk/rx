package rx

// Last emits only the last value emitted by the source [Observable].
// If the source turns out to be empty, Last emits an [Error] notification
// of [ErrEmpty].
func Last[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				var last struct {
					value    T
					hasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						last.value = n.Value
						last.hasValue = true

					case KindComplete:
						if last.hasValue {
							Try1(o, Next(last.value), func() { o.Stop(ErrOops) })
							o.Complete()
						} else {
							o.Error(ErrEmpty)
						}

					case KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}

// LastOrElse emits only the last value emitted by the source [Observable].
// If the source turns out to be empty, LastOrElse emits a specified default
// value.
func LastOrElse[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				var last struct {
					value    T
					hasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						last.value = n.Value
						last.hasValue = true

					case KindComplete:
						v := def

						if last.hasValue {
							v = last.value
						}

						Try1(o, Next(v), func() { o.Stop(ErrOops) })
						o.Emit(n)

					case KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}
