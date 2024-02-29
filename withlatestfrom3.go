package rx

// WithLatestFrom3 combines the source with 3 other Observables to create
// an Observable that emits projections of the latest values emitted by
// each Observable, only when the source emits.
func WithLatestFrom3[T0, T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v0 T0, v1 T1, v2 T2, v3 T3) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom4(source, obs1, obs2, obs3, proj)
		},
	)
}

func withLatestFrom4[T1, T2, T3, T4, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4) R,
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
		chan4 := make(chan Notification[T4])

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { obs3.Subscribe(c, channelObserver(chan3, noop)) })
		c.Go(func() { obs4.Subscribe(c, channelObserver(chan4, noop)) })

		c.Go(func() {
			var s withLatestFromState4[T1, T2, T3, T4]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = withLatestFromTry4(sink, n, proj, &s, &s.V1, 1)
				case n := <-chan2:
					cont = withLatestFromTry4(sink, n, proj, &s, &s.V2, 2)
				case n := <-chan3:
					cont = withLatestFromTry4(sink, n, proj, &s, &s.V3, 4)
				case n := <-chan4:
					cont = withLatestFromTry4(sink, n, proj, &s, &s.V4, 8)
				}
			}
		})
	}
}

type withLatestFromState4[T1, T2, T3, T4 any] struct {
	NBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

func withLatestFromTry4[T1, T2, T3, T4, R, X any](
	sink Observer[R],
	n Notification[X],
	proj func(T1, T2, T3, T4) R,
	s *withLatestFromState4[T1, T2, T3, T4],
	v *X,
	bit uint8,
) bool {
	const FullBits = 15

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			oops := func() { sink.Error(ErrOops) }
			v := Try41(proj, s.V1, s.V2, s.V3, s.V4, oops)
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
