package rx

// WithLatestFrom1 combines the source with another Observable to create
// an Observable that emits mappings of the latest values emitted by each
// Observable, only when the source emits.
func WithLatestFrom1[T0, T1, R any](
	obs1 Observable[T1],
	mapping func(v0 T0, v1 T1) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom2(source, obs1, mapping)
		},
	)
}

func withLatestFrom2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
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

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })

		c.Go(func() {
			var s withLatestFromState2[T1, T2]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = withLatestFromTry2(sink, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = withLatestFromTry2(sink, n, mapping, &s, &s.V2, 2)
				}
			}
		})
	}
}

type withLatestFromState2[T1, T2 any] struct {
	NBits uint8

	V1 T1
	V2 T2
}

func withLatestFromTry2[T1, T2, R, X any](
	sink Observer[R],
	n Notification[X],
	mapping func(T1, T2) R,
	s *withLatestFromState2[T1, T2],
	v *X,
	bit uint8,
) bool {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			oops := func() { sink.Error(ErrOops) }
			v := Try21(mapping, s.V1, s.V2, oops)
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
