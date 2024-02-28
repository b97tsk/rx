package rx

// SkipAll skips all values emitted by the source Observable.
func SkipAll[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindError, KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}
