package rx

// StartWith emits the items you specify as arguments before it begins to
// mirrors the source Observable.
func StartWith[T any](s ...T) Operator[T, T] {
	if len(s) == 0 {
		return AsOperator(identity[Observable[T]])
	}

	return startWith(s...)
}

func startWith[T any](s ...T) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return Concat(FromSlice(s), source)
		},
	)
}
