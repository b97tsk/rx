package rx

// SkipLast skips the last count values emitted by the source Observable.
func SkipLast[T any](count int) Operator[T, T] {
	if count <= 0 {
		return NewOperator(identity[Observable[T]])
	}

	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				b := make([]T, 0, count)
				i := 0

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if len(b) < cap(b) {
							b = b[:len(b)+1]
						} else {
							sink.Next(b[i])
						}

						b[i] = n.Value
						i = (i + 1) % cap(b)

					case KindError, KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}
