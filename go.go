package rx

// Go mirrors the source Observable in a goroutine.
func Go[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c.Go(func() { source.Subscribe(c, sink) })
			}
		},
	)
}
