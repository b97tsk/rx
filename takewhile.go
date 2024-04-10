package rx

// TakeWhile emits values emitted by the source Observable so long as
// each value satisfies a given predicate function, and then completes
// as soon as the predicate function returns false.
func TakeWhile[T any](pred func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						if pred(n.Value) {
							o.Emit(n)
							return
						}

						noop = true
						o.Complete()

					case KindError, KindComplete:
						o.Emit(n)
					}
				})
			}
		},
	)
}
