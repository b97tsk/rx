package rx

// Contains emits a boolean to indicate whether any value of the source
// [Observable] satisfies a given predicate function.
func Contains[T any](pred func(v T) bool) Operator[T, bool] {
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
						if pred(n.Value) {
							o.Next(true)
							noop = true
							o.Complete()
						}
					case KindComplete:
						Try1(o, Next(false), func() { o.Stop(ErrOops) })
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

// ContainsElement emits a boolean to indicate whether the source [Observable]
// emits a given value.
func ContainsElement[T comparable](v T) Operator[T, bool] {
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
						if n.Value == v {
							o.Next(true)
							noop = true
							o.Complete()
						}
					case KindComplete:
						Try1(o, Next(false), func() { o.Stop(ErrOops) })
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
