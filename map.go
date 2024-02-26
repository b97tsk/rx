package rx

// Map applies a given projection function to each value emitted by
// the source Observable, then emits the resulting values.
func Map[T, R any](proj func(v T) R) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return map1(proj)
}

func map1[T, R any](proj func(v T) R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, sink Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						sink.Next(proj(n.Value))
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

// MapTo emits the given constant value on the output Observable every time
// the source Observable emits a value.
func MapTo[T, R any](v R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, sink Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						sink.Next(v)
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
