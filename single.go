package rx

// Single emits the single value emitted by the source [Observable].
// If the source emits more than one value or no values, it emits
// an [Error] notification of [ErrNotSingle] or [ErrEmpty] respectively.
func Single[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var x struct {
					first struct {
						value    T
						hasValue bool
					}
					noop bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					if x.noop {
						return
					}

					switch n.Kind {
					case KindNext:
						if !x.first.hasValue {
							x.first.value = n.Value
							x.first.hasValue = true
							return
						}

						x.noop = true
						o.Error(ErrNotSingle)

					case KindComplete:
						if x.first.hasValue {
							Try1(o, Next(x.first.value), func() { o.Stop(ErrOops) })
							o.Complete()
						} else {
							o.Error(ErrEmpty)
						}

					case KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}
