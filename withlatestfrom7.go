package rx

// WithLatestFrom7 combines the source with 7 other Observables to create
// an Observable that emits projection of latest values of each Observable,
// only when the source emits.
func WithLatestFrom7[T0, T1, T2, T3, T4, T5, T6, T7, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	proj func(v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom8(source, obs1, obs2, obs3, obs4, obs5, obs6, obs7, proj)
		},
	)
}

func withLatestFrom8[T1, T2, T3, T4, T5, T6, T7, T8, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	obs8 Observable[T8],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8) R,
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
		chan5 := make(chan Notification[T5])
		chan6 := make(chan Notification[T6])
		chan7 := make(chan Notification[T7])
		chan8 := make(chan Notification[T8])

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { obs3.Subscribe(c, channelObserver(chan3, noop)) })
		c.Go(func() { obs4.Subscribe(c, channelObserver(chan4, noop)) })
		c.Go(func() { obs5.Subscribe(c, channelObserver(chan5, noop)) })
		c.Go(func() { obs6.Subscribe(c, channelObserver(chan6, noop)) })
		c.Go(func() { obs7.Subscribe(c, channelObserver(chan7, noop)) })
		c.Go(func() { obs8.Subscribe(c, channelObserver(chan8, noop)) })

		c.Go(func() {
			var s withLatestFromState8[T1, T2, T3, T4, T5, T6, T7, T8]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V1, 1)
				case n := <-chan2:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V2, 2)
				case n := <-chan3:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V3, 4)
				case n := <-chan4:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V4, 8)
				case n := <-chan5:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V5, 16)
				case n := <-chan6:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V6, 32)
				case n := <-chan7:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V7, 64)
				case n := <-chan8:
					cont = withLatestFromTry8(sink, n, proj, &s, &s.V8, 128)
				}
			}
		})
	}
}

type withLatestFromState8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	NBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
	V8 T8
}

func withLatestFromTry8[T1, T2, T3, T4, T5, T6, T7, T8, R, X any](
	sink Observer[R],
	n Notification[X],
	proj func(T1, T2, T3, T4, T5, T6, T7, T8) R,
	s *withLatestFromState8[T1, T2, T3, T4, T5, T6, T7, T8],
	v *X,
	bit uint8,
) bool {
	const FullBits = 255

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			oops := func() { sink.Error(ErrOops) }
			v := Try81(proj, s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7, s.V8, oops)
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
