package rx

// SkipLast skips the last count values emitted by the source [Observable].
func SkipLast[T any](count int) Operator[T, T] {
	if count <= 0 {
		return NewOperator(identity[Observable[T]])
	}

	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				buf := make([]T, 0, count)
				i := 0

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if len(buf) < cap(buf) {
							buf = buf[:len(buf)+1]
						} else {
							o.Next(buf[i])
						}

						buf[i] = n.Value
						i = (i + 1) % cap(buf)

					case KindComplete, KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}
