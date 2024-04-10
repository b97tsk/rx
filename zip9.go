package rx

// Zip9 combines multiple Observables to create an Observable that emits
// mappings of the values emitted by each of its input Observables.
//
// Zip9 pulls values from each input Observable one by one, it does not
// buffer any value.
func Zip9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	obs8 Observable[T8],
	obs9 Observable[T9],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9) R,
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
		chan6 := make(chan Notification[T6])
		chan7 := make(chan Notification[T7])
		chan8 := make(chan Notification[T8])
		chan9 := make(chan Notification[T9])

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
			Again6:
				n6 := <-chan6
				switch n6.Kind {
				case KindNext:
				case KindError:
					o.Error(n6.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again6
				}
			Again7:
				n7 := <-chan7
				switch n7.Kind {
				case KindNext:
				case KindError:
					o.Error(n7.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again7
				}
			Again8:
				n8 := <-chan8
				switch n8.Kind {
				case KindNext:
				case KindError:
					o.Error(n8.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again8
				}
			Again9:
				n9 := <-chan9
				switch n9.Kind {
				case KindNext:
				case KindError:
					o.Error(n9.Error)
					return
				case KindComplete:
					o.Complete()
					return
				default:
					goto Again9
				}
				v := Try91(mapping, n1.Value, n2.Value, n3.Value, n4.Value, n5.Value, n6.Value, n7.Value, n8.Value, n9.Value, oops)
				Try1(o, Next(v), oops)
			}
		})

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { obs3.Subscribe(c, channelObserver(chan3, noop)) })
		c.Go(func() { obs4.Subscribe(c, channelObserver(chan4, noop)) })
		c.Go(func() { obs5.Subscribe(c, channelObserver(chan5, noop)) })
		c.Go(func() { obs6.Subscribe(c, channelObserver(chan6, noop)) })
		c.Go(func() { obs7.Subscribe(c, channelObserver(chan7, noop)) })
		c.Go(func() { obs8.Subscribe(c, channelObserver(chan8, noop)) })
		c.Go(func() { obs9.Subscribe(c, channelObserver(chan9, noop)) })
	}
}
