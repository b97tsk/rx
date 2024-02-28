package rx

// Filter filters values emitted by the source Observable by only emitting
// those that satisfy a given condition.
func Filter[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if cond(n.Value) {
							sink(n)
						}
					case KindError, KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}

// FilterOut filters out values emitted by the source Observable by only
// emitting those that do not satisfy a given condition.
func FilterOut[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if !cond(n.Value) {
							sink(n)
						}
					case KindError, KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}

// FilterMap passes each value emitted by the source Observable to a given
// condition function and emits their mapping, the first return value of
// the condition function, only if the second is true.
func FilterMap[T, R any](cond func(v T) (R, bool)) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, sink Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if v, ok := cond(n.Value); ok {
							sink.Next(v)
						}
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						sink.Complete()
					}
				})
			}
		},
	)
}
