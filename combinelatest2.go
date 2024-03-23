package rx

// CombineLatest2 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest2[T1, T2, R any](
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
			var s combineLatestState2[T1, T2]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestTry2(sink, n, mapping, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestTry2(sink, n, mapping, &s, &s.V2, 2)
				}
			}
		})
	}
}

type combineLatestState2[T1, T2 any] struct {
	NBits, CBits uint8

	V1 T1
	V2 T2
}

func combineLatestTry2[T1, T2, R, X any](
	sink Observer[R],
	n Notification[X],
	mapping func(T1, T2) R,
	s *combineLatestState2[T1, T2],
	v *X,
	bit uint8,
) bool {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			oops := func() { sink.Error(ErrOops) }
			v := Try21(mapping, s.V1, s.V2, oops)
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
