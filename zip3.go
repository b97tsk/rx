package rx

// Zip3 combines multiple Observables to create an [Observable] that emits
// mappings of the values emitted by each of the input Observables.
//
// Zip3 pulls values from each input [Observable] one by one, it does not
// buffer any value.
func Zip3[T1, T2, T3, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	mapping func(v1 T1, v2 T2, v3 T3) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, cancel := c.WithCancel()
		noop := make(chan struct{})
		o = o.DoOnTermination(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])

		c.Go(func() {
			oops := func() { o.Stop(ErrOops) }
			for {
			Again1:
				n1 := <-chan1
				switch n1.Kind {
				case KindNext:
				case KindComplete:
					o.Complete()
					return
				case KindError:
					o.Error(n1.Error)
					return
				case KindStop:
					o.Stop(n1.Error)
					return
				default:
					goto Again1
				}
			Again2:
				n2 := <-chan2
				switch n2.Kind {
				case KindNext:
				case KindComplete:
					o.Complete()
					return
				case KindError:
					o.Error(n2.Error)
					return
				case KindStop:
					o.Stop(n2.Error)
					return
				default:
					goto Again2
				}
			Again3:
				n3 := <-chan3
				switch n3.Kind {
				case KindNext:
				case KindComplete:
					o.Complete()
					return
				case KindError:
					o.Error(n3.Error)
					return
				case KindStop:
					o.Stop(n3.Error)
					return
				default:
					goto Again3
				}
				v := Try31(mapping, n1.Value, n2.Value, n3.Value, oops)
				Try1(o, Next(v), oops)
			}
		})

		c.Go(func() { ob1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { ob2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { ob3.Subscribe(c, channelObserver(chan3, noop)) })
	}
}
