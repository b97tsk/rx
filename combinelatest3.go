package rx

// CombineLatest3 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest3[T1, T2, T3, R any](
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
			var s combineLatestState3[T1, T2, T3]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestEmit3(o, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestEmit3(o, n, mapping, &s, &s.V2, 2)
				case n := <-chan3:
					cont = combineLatestEmit3(o, n, mapping, &s, &s.V3, 4)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, ob1, chan1, noop) &&
			subscribeChannel(c, ob2, chan2, noop) &&
			subscribeChannel(c, ob3, chan3, noop)
	}
}

type combineLatestState3[T1, T2, T3 any] struct {
	NBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
}

func combineLatestEmit3[T1, T2, T3, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3) R,
	s *combineLatestState3[T1, T2, T3],
	v *X,
	bit uint8,
) bool {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			oops := func() { o.Error(ErrOops) }
			v := Try31(mapping, s.V1, s.V2, s.V3, oops)
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
