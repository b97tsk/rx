package rx

// Zip5 combines multiple Observables to create an Observable that emits
// mappings of the values emitted by each of its input Observables.
//
// Zip5 pulls values from each input Observable one by one, it does not
// buffer any value.
func Zip5[T1, T2, T3, T4, T5, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
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
		chan4 := make(chan Notification[T4])
		chan5 := make(chan Notification[T5])

		c.Go(func() {
			oops := func() { o.Error(ErrOops) }
			for {
			Again1:
				n1 := <-chan1
				switch n1.Kind {
				case KindNext:
				case KindError:
					o.Error(n1.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again1
				}
			Again2:
				n2 := <-chan2
				switch n2.Kind {
				case KindNext:
				case KindError:
					o.Error(n2.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again2
				}
			Again3:
				n3 := <-chan3
				switch n3.Kind {
				case KindNext:
				case KindError:
					o.Error(n3.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again3
				}
			Again4:
				n4 := <-chan4
				switch n4.Kind {
				case KindNext:
				case KindError:
					o.Error(n4.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again4
				}
			Again5:
				n5 := <-chan5
				switch n5.Kind {
				case KindNext:
				case KindError:
					o.Error(n5.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again5
				}
				v := Try51(mapping, n1.Value, n2.Value, n3.Value, n4.Value, n5.Value, oops)
				Try1(o, Next(v), oops)
			}
		})

		c.Go(func() { ob1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { ob2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { ob3.Subscribe(c, channelObserver(chan3, noop)) })
		c.Go(func() { ob4.Subscribe(c, channelObserver(chan4, noop)) })
		c.Go(func() { ob5.Subscribe(c, channelObserver(chan5, noop)) })
	}
}
