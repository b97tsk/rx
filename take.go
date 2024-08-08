package rx

// Take emits only the first count values emitted by the source [Observable].
func Take[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count <= 0 {
				return Empty[T]()
			}

			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var x struct {
					count int
					noop  bool
				}

				x.count = count

				source.Subscribe(c, func(n Notification[T]) {
					if x.noop {
						return
					}

					o.Emit(n)

					if n.Kind == KindNext {
						x.count--
						if x.count == 0 {
							x.noop = true
							o.Complete()
						}
					}
				})
			}
		},
	)
}
