package rx

// Contains emits a boolean to indicate whether any value of the source
// Observable satisfies a given predicate function.
func Contains[T any](pred func(v T) bool) Operator[T, bool] {
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
						if pred(n.Value) {
							sink.Next(true)
							noop = true
							sink.Complete()
						}
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						Try1(sink, Next(false), func() { sink.Error(ErrOops) })
						sink.Complete()
					}
				})
			}
		},
	)
}

// ContainsElement emits a boolean to indicate whether the source Observable
// emits a given value.
func ContainsElement[T comparable](v T) Operator[T, bool] {
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
						if n.Value == v {
							sink.Next(true)
							noop = true
							sink.Complete()
						}
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						Try1(sink, Next(false), func() { sink.Error(ErrOops) })
						sink.Complete()
					}
				})
			}
		},
	)
}
