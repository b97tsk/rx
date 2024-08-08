package rx

// Channelize splits upstream and downstream with two channels, then uses
// a given join function to connect them.
//
// Channelize closes downstream channel after join returns.
func Channelize[T any](join func(upstream <-chan Notification[T], downstream chan<- Notification[T])) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				c, cancel := c.WithCancel()
				o = o.DoOnTermination(cancel)

				noop := make(chan struct{})
				upstream := make(chan Notification[T])
				downstream := make(chan Notification[T])

				c.Go(func() {
					defer func() {
						cancel()
						close(noop)
						close(downstream)
					}()
					Try2(join, upstream, downstream, func() { downstream <- Stop[T](ErrOops) })
				})

				c.Go(func() {
					for n := range downstream {
						switch n.Kind {
						case KindNext:
							Try1(o, n, func() {
								c.Go(func() { drain(downstream) })
								o.Stop(ErrOops)
							})
						case KindComplete, KindError, KindStop:
							c.Go(func() { drain(downstream) })
							o.Emit(n)
							return
						}
					}

					defer o.Stop(ErrOops)
					panic("Channelize: no termination")
				})

				source.Subscribe(c, channelObserver(upstream, noop))
			}
		},
	)
}

func drain[T any](c <-chan T) {
	for range c {
	}
}
