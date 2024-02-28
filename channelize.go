package rx

// Channelize separates upstream and downstream with two channels, then uses
// provided join function to connect them.
//
// Notifications sent to downstream must honor the [Observable] protocol.
//
// Channelize closes downstream channel after join returns.
func Channelize[T any](join func(upstream <-chan Notification[T], downstream chan<- Notification[T])) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.OnLastNotification(cancel)

				upstream := make(chan Notification[T])
				downstream := make(chan Notification[T])
				noop := make(chan struct{})

				c.Go(func() {
					defer func() {
						close(downstream)
						close(noop)
					}()
					Try2(join, upstream, downstream, func() { downstream <- Error[T](ErrOops) })
				})

				c.Go(func() {
					for n := range downstream {
						switch n.Kind {
						case KindNext:
							Try1(sink, n, func() {
								c.Go(func() { drain(downstream) })
								sink.Error(ErrOops)
							})
						case KindError, KindComplete:
							defer drain(downstream)
							sink(n)
							return
						}
					}
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
