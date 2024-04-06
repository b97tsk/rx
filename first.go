package rx

// First emits only the first value emitted by the source Observable.
// If the source turns out to be empty, First emits a notification of ErrEmpty.
func First[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						sink(n)
						noop = true
						sink.Complete()
					case KindError:
						sink(n)
					case KindComplete:
						sink.Error(ErrEmpty)
					}
				})
			}
		},
	)
}

// FirstOrElse emits only the first value emitted by the source Observable.
// If the source turns out to be empty, FirstOrElse emits a specified default
// value.
func FirstOrElse[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						sink(n)
						noop = true
						sink.Complete()
					case KindError:
						sink(n)
					case KindComplete:
						Try1(sink, Next(def), func() { sink.Error(ErrOops) })
						sink(n)
					}
				})
			}
		},
	)
}
