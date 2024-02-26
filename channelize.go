package rx

// Channelize separates upstream and downstream with two channels, then uses
// provided join function to connect them.
//
// Notifications sent to downstream must honor the [Observable] protocol.
//
// Channelize closes downstream channel after join returns.
func Channelize[T any](join func(upstream <-chan Notification[T], downstream chan<- Notification[T])) Operator[T, T] {
	if join == nil {
		panic("join == nil")
	}

	return channelize[T](join)
}

func channelize[T any](join func(upstream <-chan Notification[T], downstream chan<- Notification[T])) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				c, cancel := c.WithCancel()
				sink = sink.OnLastNotification(cancel)

				upstream := make(chan Notification[T])
				downstream := make(chan Notification[T])
				noop := make(chan struct{})

				c.Go(func() {
					join(upstream, downstream)
					close(downstream)
					close(noop)
				})

				c.Go(func() {
					for n := range downstream {
						sink(n)
					}
				})

				source.Subscribe(c, channelObserver(upstream, noop))
			}
		},
	)
}
