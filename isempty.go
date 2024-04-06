package rx

// IsEmpty emits a boolean to indicate whether the source emits no values.
func IsEmpty[T any]() Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return func(c Context, sink Observer[bool]) {
				c, cancel := c.WithCancel()
				sink = sink.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						sink.Next(false)
						noop = true
						sink.Complete()
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
