package rx

// CombineLatest5 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest5[T1, T2, T3, T4, T5, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
) Observable[R] {
	return func(c Context, sink Observer[R]) {
		c, cancel := c.WithCancel()
		noop := make(chan struct{})
		sink = sink.OnTermination(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])
		chan4 := make(chan Notification[T4])
		chan5 := make(chan Notification[T5])

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { obs3.Subscribe(c, channelObserver(chan3, noop)) })
		c.Go(func() { obs4.Subscribe(c, channelObserver(chan4, noop)) })
		c.Go(func() { obs5.Subscribe(c, channelObserver(chan5, noop)) })

		c.Go(func() {
			var s combineLatestState5[T1, T2, T3, T4, T5]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestTry5(sink, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestTry5(sink, n, mapping, &s, &s.V2, 2)
				case n := <-chan3:
					cont = combineLatestTry5(sink, n, mapping, &s, &s.V3, 4)
				case n := <-chan4:
					cont = combineLatestTry5(sink, n, mapping, &s, &s.V4, 8)
				case n := <-chan5:
					cont = combineLatestTry5(sink, n, mapping, &s, &s.V5, 16)
				}
			}
		})
	}
}

type combineLatestState5[T1, T2, T3, T4, T5 any] struct {
	NBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
}

func combineLatestTry5[T1, T2, T3, T4, T5, R, X any](
	sink Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5) R,
	s *combineLatestState5[T1, T2, T3, T4, T5],
	v *X,
	bit uint8,
) bool {
	const FullBits = 31

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			oops := func() { sink.Error(ErrOops) }
			v := Try51(mapping, s.V1, s.V2, s.V3, s.V4, s.V5, oops)
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
