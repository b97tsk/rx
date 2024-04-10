package rx

// Single emits the single value emitted by the source Observable.
// If the source emits more than one value or no values, it emits
// a notification of ErrNotSingle or ErrEmpty respectively.
func Single[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

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
						o.Error(ErrNotSingle)

					case KindError:
						o.Emit(n)

					case KindComplete:
						if first.HasValue {
							Try1(o, Next(first.Value), func() { o.Error(ErrOops) })
							o.Complete()
						} else {
							o.Error(ErrEmpty)
						}
					}
				})
			}
		},
	)
}
