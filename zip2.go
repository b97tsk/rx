package rx

// Zip2 combines multiple Observables to create an Observable that emits
// mappings of the values emitted by each of its input Observables.
//
// Zip2 pulls values from each input Observable one by one, it does not
// buffer any value.
func Zip2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, sink Observer[R]) {
		c, cancel := c.WithCancel()
		noop := make(chan struct{})
		sink = sink.DoOnTermination(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		c.Go(func() {
			oops := func() { sink.Error(ErrOops) }
			for {
			Again1:
				n1 := <-chan1
				switch n1.Kind {
				case KindNext:
				case KindError:
					sink.Error(n1.Error)
					return
				case KindComplete:
					sink.Complete()
					return
				default:
					goto Again1
				}
			Again2:
				n2 := <-chan2
				switch n2.Kind {
				case KindNext:
				case KindError:
					sink.Error(n2.Error)
					return
				case KindComplete:
					sink.Complete()
					return
				default:
					goto Again2
				}
				v := Try21(mapping, n1.Value, n2.Value, oops)
				Try1(sink, Next(v), oops)
			}
		})

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
	}
}
