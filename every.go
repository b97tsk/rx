package rx

// Every emits a boolean to indicate whether every value of the source
// Observable satisfies a given predicate function.
func Every[T any](pred func(v T) bool) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return func(c Context, sink Observer[bool]) {
				c, cancel := c.WithCancel()
				sink = sink.OnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						if !pred(n.Value) {
							sink.Next(false)
							noop = true
							sink.Complete()
						}
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						Try1(sink, Next(true), func() { sink.Error(ErrOops) })
						sink.Complete()
					}
				})
			}
		},
	)
}
