package rx

import "context"

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
			return func(ctx context.Context, sink Observer[T]) {
				wg := WaitGroupFromContext(ctx)
				ctx, cancel := context.WithCancel(ctx)
				sink = sink.OnLastNotification(cancel)

				upstream := make(chan Notification[T])
				downstream := make(chan Notification[T])
				noop := make(chan struct{})

				wg.Go(func() {
					join(upstream, downstream)
					close(downstream)
					close(noop)
				})

				wg.Go(func() {
					for n := range downstream {
						sink(n)
					}
				})

				source.Subscribe(ctx, chanObserver(upstream, noop))
			}
		},
	)
}
