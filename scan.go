package rx

// Scan applies an accumulator function over the source Observable,
// and emits each intermediate result, given an initial value.
func Scan[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, sink Observer[R]) {
				res := init

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						res = accumulator(res, n.Value)
						sink.Next(res)
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
