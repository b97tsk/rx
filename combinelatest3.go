package rx

// CombineLatest3 combines multiple Observables to create an Observable
// that emits projections of the latest values emitted by each of its
// input Observables.
func CombineLatest3[T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v1 T1, v2 T2, v3 T3) R,
) Observable[R] {
	return func(c Context, sink Observer[R]) {
		c, cancel := c.WithCancel()
		noop := make(chan struct{})
		sink = sink.OnLastNotification(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { obs3.Subscribe(c, channelObserver(chan3, noop)) })

		c.Go(func() {
			var s combineLatestState3[T1, T2, T3]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestTry3(sink, n, proj, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestTry3(sink, n, proj, &s, &s.V2, 2)
				case n := <-chan3:
					cont = combineLatestTry3(sink, n, proj, &s, &s.V3, 4)
				}
			}
		})
	}
}

type combineLatestState3[T1, T2, T3 any] struct {
	NBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
}

func combineLatestTry3[T1, T2, T3, R, X any](
	sink Observer[R],
	n Notification[X],
	proj func(T1, T2, T3) R,
	s *combineLatestState3[T1, T2, T3],
	v *X,
	bit uint8,
) bool {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			oops := func() { sink.Error(ErrOops) }
			v := Try31(proj, s.V1, s.V2, s.V3, oops)
			Try1(sink, Next(v), oops)
		}

	case KindError:
		sink.Error(n.Error)
		return false

	case KindComplete:
		if s.CBits |= bit; s.CBits == FullBits {
			sink.Complete()
			return false
		}
	}

	return true
}
