package rx

// Take emits only the first count values emitted by the source Observable.
func Take[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count <= 0 {
				return Empty[T]()
			}

			return takeObservable[T]{source, count}.Subscribe
		},
	)
}

type takeObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (obs takeObservable[T]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var noop bool

	count := obs.Count

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		sink(n)

		if n.Kind == KindNext {
			count--

			if count == 0 {
				noop = true
				sink.Complete()
			}
		}
	})
}
