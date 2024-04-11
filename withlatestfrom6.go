package rx

// WithLatestFrom6 combines the source with 6 other Observables to create
// an Observable that emits mappings of the latest values emitted by each
// Observable, only when the source emits.
func WithLatestFrom6[T0, T1, T2, T3, T4, T5, T6, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	ob6 Observable[T6],
	mapping func(v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom7(source, ob1, ob2, ob3, ob4, ob5, ob6, mapping)
		},
	)
}

func withLatestFrom7[T1, T2, T3, T4, T5, T6, T7, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	ob6 Observable[T6],
	ob7 Observable[T7],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
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

		c.Go(func() {
			var s withLatestFromState7[T1, T2, T3, T4, T5, T6, T7]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V2, 2)
				case n := <-chan3:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V3, 4)
				case n := <-chan4:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V4, 8)
				case n := <-chan5:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V5, 16)
				case n := <-chan6:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V6, 32)
				case n := <-chan7:
					cont = withLatestFromTry7(o, n, mapping, &s, &s.V7, 64)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, ob1, chan1, noop) &&
			subscribeChannel(c, ob2, chan2, noop) &&
			subscribeChannel(c, ob3, chan3, noop) &&
			subscribeChannel(c, ob4, chan4, noop) &&
			subscribeChannel(c, ob5, chan5, noop) &&
			subscribeChannel(c, ob6, chan6, noop) &&
			subscribeChannel(c, ob7, chan7, noop)
	}
}

type withLatestFromState7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	NBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
}

func withLatestFromTry7[T1, T2, T3, T4, T5, T6, T7, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5, T6, T7) R,
	s *withLatestFromState7[T1, T2, T3, T4, T5, T6, T7],
	v *X,
	bit uint8,
) bool {
	const FullBits = 127

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			oops := func() { o.Error(ErrOops) }
			v := Try71(mapping, s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7, oops)
			Try1(o, Next(v), oops)
		}

	case KindError:
		o.Error(n.Error)
		return false

	case KindComplete:
		if bit == 1 {
			o.Complete()
			return false
		}
	}

	return true
}
