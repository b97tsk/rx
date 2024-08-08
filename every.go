package rx

// Every emits a boolean to indicate whether every value of the source
// [Observable] satisfies a given predicate function.
func Every[T any](pred func(v T) bool) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return func(c Context, o Observer[bool]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				var noop bool

				source.Subscribe(c, func(n Notification[T]) {
					if noop {
						return
					}

					switch n.Kind {
					case KindNext:
						if !pred(n.Value) {
							o.Next(false)
							noop = true
							o.Complete()
						}
					case KindComplete:
						Try1(o, Next(true), func() { o.Stop(ErrOops) })
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
