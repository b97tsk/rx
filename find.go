package rx

// Find emits only the first value emitted by the source Observable that
// satisfies a given predicate function.
func Find[T any](pred func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.OnLastNotification(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						if pred(n.Value) {
							sink(n)
							noop = true
							sink.Complete()
						}
					case KindError, KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}
