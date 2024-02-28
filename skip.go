package rx

// Skip skips the first count values emitted by the source Observable.
func Skip[T any](count int) Operator[T, T] {
	if count <= 0 {
		return NewOperator(identity[Observable[T]])
	}

	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				var taking bool

				count := count

				source.Subscribe(c, func(n Notification[T]) {
					if taking || n.Kind != KindNext {
						sink(n)
						return
					}

					count--
					taking = count == 0
				})
			}
		},
	)
}
