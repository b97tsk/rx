package rx

// Last emits only the last value emitted by the source Observable.
// If the source turns out to be empty, Last emits a notification of ErrEmpty.
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
					case KindError:
						o.Emit(n)
					case KindComplete:
						if last.hasValue {
							Try1(o, Next(last.value), func() { o.Error(ErrOops) })
							o.Complete()
						} else {
							o.Error(ErrEmpty)
						}
					}
				})
			}
		},
	)
}

// LastOrElse emits only the last value emitted by the source Observable.
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
					case KindError:
						o.Emit(n)
					case KindComplete:
						v := def

						if last.hasValue {
							v = last.value
						}

						Try1(o, Next(v), func() { o.Error(ErrOops) })
						o.Emit(n)
					}
				})
			}
		},
	)
}
