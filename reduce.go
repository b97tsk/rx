package rx

// Reduce applies an accumulator function over the source Observable,
// and emits the accumulated result when the source completes, given
// an initial value.
func Reduce[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	if accumulator == nil {
		panic("accumulator == nil")
	}

	return reduce(init, accumulator)
}

func reduce[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, sink Observer[R]) {
				res := init

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						res = accumulator(res, n.Value)
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						sink.Next(res)
						sink.Complete()
					}
				})
			}
		},
	)
}
