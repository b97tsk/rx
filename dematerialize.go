package rx

// Dematerialize converts an [Observable] of [Notification] values into
// the emissions that they represent. It's the opposite of [Materialize].
func Dematerialize[_ Notification[T], T any]() Operator[Notification[T], T] {
	return NewOperator(
		func(source Observable[Notification[T]]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[Notification[T]]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						n := n.Value

						switch n.Kind {
						case KindNext:
						case KindComplete, KindError, KindStop:
							noop = true
						}

						o.Emit(n)
					case KindComplete:
						o.Complete()
					case KindError:
						o.Error(n.Error)
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}
