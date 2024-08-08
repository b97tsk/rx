package rx

// Filter filters values emitted by the source [Observable] by only emitting
// those that satisfy a given predicate function.
func Filter[T any](pred func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if pred(n.Value) {
							o.Emit(n)
						}
					case KindComplete, KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}

// FilterOut filters out values emitted by the source [Observable] by only
// emitting those that do not satisfy a given predicate function.
func FilterOut[T any](pred func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if !pred(n.Value) {
							o.Emit(n)
						}
					case KindComplete, KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}

// FilterMap passes each value emitted by the source [Observable] to a given
// predicate function and emits their mapping, the first return value of
// the predicate function, only if the second is true.
func FilterMap[T, R any](pred func(v T) (R, bool)) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, o Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if v, ok := pred(n.Value); ok {
							o.Next(v)
						}
					case KindComplete:
						o.Complete()
					case KindError:
						o.Error(n.Error)
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}
