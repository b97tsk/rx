package rx

// CombineLatest4 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest4[T1, T2, T3, T4, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4) R,
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

		c.Go(func() {
			var s combineLatestState4[T1, T2, T3, T4]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestTry4(o, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestTry4(o, n, mapping, &s, &s.V2, 2)
				case n := <-chan3:
					cont = combineLatestTry4(o, n, mapping, &s, &s.V3, 4)
				case n := <-chan4:
					cont = combineLatestTry4(o, n, mapping, &s, &s.V4, 8)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, obs1, chan1, noop) &&
			subscribeChannel(c, obs2, chan2, noop) &&
			subscribeChannel(c, obs3, chan3, noop) &&
			subscribeChannel(c, obs4, chan4, noop)
	}
}

type combineLatestState4[T1, T2, T3, T4 any] struct {
	NBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

func combineLatestTry4[T1, T2, T3, T4, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4) R,
	s *combineLatestState4[T1, T2, T3, T4],
	v *X,
	bit uint8,
) bool {
	const FullBits = 15

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			oops := func() { o.Error(ErrOops) }
			v := Try41(mapping, s.V1, s.V2, s.V3, s.V4, oops)
			Try1(o, Next(v), oops)
		}

	case KindError:
		o.Error(n.Error)
		return false

	case KindComplete:
		if s.CBits |= bit; s.CBits == FullBits {
			o.Complete()
			return false
		}
	}

	return true
}
