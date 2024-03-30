package rx

// Single emits the single value emitted by the source Observable.
// If the source emits more than one value or no values, it emits
// a notification of ErrNotSingle or ErrEmpty respectively.
func Single[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.DoOnTermination(cancel)

				var first struct {
					Value    T
					HasValue bool
				}

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						if !first.HasValue {
							first.Value = n.Value
							first.HasValue = true
							return
						}

						noop = true
						sink.Error(ErrNotSingle)

					case KindError:
						sink(n)

					case KindComplete:
						if first.HasValue {
							Try1(sink, Next(first.Value), func() { sink.Error(ErrOops) })
							sink.Complete()
						} else {
							sink.Error(ErrEmpty)
						}
					}
				})
			}
		},
	)
}
