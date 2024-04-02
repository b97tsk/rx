package rx

// WithLatestFrom2 combines the source with 2 other Observables to create
// an Observable that emits mappings of the latest values emitted by each
// Observable, only when the source emits.
func WithLatestFrom2[T0, T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	mapping func(v0 T0, v1 T1, v2 T2) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom3(source, obs1, obs2, mapping)
		},
	)
}

func withLatestFrom3[T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	mapping func(v1 T1, v2 T2, v3 T3) R,
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
		chan3 := make(chan Notification[T3])

		c.Go(func() {
			var s withLatestFromState3[T1, T2, T3]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = withLatestFromTry3(sink, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = withLatestFromTry3(sink, n, mapping, &s, &s.V2, 2)
				case n := <-chan3:
					cont = withLatestFromTry3(sink, n, mapping, &s, &s.V3, 4)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, obs1, chan1, noop) &&
			subscribeChannel(c, obs2, chan2, noop) &&
			subscribeChannel(c, obs3, chan3, noop)
	}
}

type withLatestFromState3[T1, T2, T3 any] struct {
	NBits uint8

	V1 T1
	V2 T2
	V3 T3
}

func withLatestFromTry3[T1, T2, T3, R, X any](
	sink Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3) R,
	s *withLatestFromState3[T1, T2, T3],
	v *X,
	bit uint8,
) bool {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			oops := func() { sink.Error(ErrOops) }
			v := Try31(mapping, s.V1, s.V2, s.V3, oops)
			Try1(sink, Next(v), oops)
		}

	case KindError:
		sink.Error(n.Error)
		return false

	case KindComplete:
		if bit == 1 {
			sink.Complete()
			return false
		}
	}

	return true
}
