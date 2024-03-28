package rx

// Take emits only the first count values emitted by the source Observable.
func Take[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count <= 0 {
				return Empty[T]()
			}

			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.OnTermination(cancel)

				var noop bool

				count := count

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					sink(n)

					if n.Kind == KindNext {
						count--

						if count == 0 {
							noop = true
							sink.Complete()
						}
					}
				})
			}
		},
	)
}
