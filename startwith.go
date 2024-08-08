package rx

// StartWith emits the values you specify as arguments before it begins to
// mirrors the source [Observable].
func StartWith[T any](s ...T) Operator[T, T] {
	if len(s) == 0 {
		return NewOperator(identity[Observable[T]])
	}

	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return Concat(FromSlice(s), source)
		},
	)
}
