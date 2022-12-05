package rx

// EndWith mirrors the source Observable, emits the items you specify as
// arguments when the source completes.
func EndWith[T any](s ...T) Operator[T, T] {
	if len(s) == 0 {
		return NewOperator(identity[Observable[T]])
	}

	return endWith(s...)
}

func endWith[T any](s ...T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return Concat(source, FromSlice(s))
		},
	)
}
