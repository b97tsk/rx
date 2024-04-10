package rx

// Map applies a given mapping function to each value emitted by the source
// Observable, then emits the resulting values.
func Map[T, R any](mapping func(v T) R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, o Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						o.Next(mapping(n.Value))
					case KindError:
						o.Error(n.Error)
					case KindComplete:
						o.Complete()
					}
				})
			}
		},
	)
}

// MapTo emits the given constant value on the output Observable every time
// the source Observable emits a value.
func MapTo[T, R any](v R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, o Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						o.Next(v)
					case KindError:
						o.Error(n.Error)
					case KindComplete:
						o.Complete()
					}
				})
			}
		},
	)
}
