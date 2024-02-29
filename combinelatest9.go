package rx

// CombineLatest9 combines multiple Observables to create an Observable
// that emits projections of the latest values emitted by each of its
// input Observables.
func CombineLatest9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	obs8 Observable[T8],
	obs9 Observable[T9],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9) R,
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
		chan9 := make(chan Notification[T9])

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })
		c.Go(func() { obs3.Subscribe(c, channelObserver(chan3, noop)) })
		c.Go(func() { obs4.Subscribe(c, channelObserver(chan4, noop)) })
		c.Go(func() { obs5.Subscribe(c, channelObserver(chan5, noop)) })
		c.Go(func() { obs6.Subscribe(c, channelObserver(chan6, noop)) })
		c.Go(func() { obs7.Subscribe(c, channelObserver(chan7, noop)) })
		c.Go(func() { obs8.Subscribe(c, channelObserver(chan8, noop)) })
		c.Go(func() { obs9.Subscribe(c, channelObserver(chan9, noop)) })

		c.Go(func() {
			var s combineLatestState9[T1, T2, T3, T4, T5, T6, T7, T8, T9]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V2, 2)
				case n := <-chan3:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V3, 4)
				case n := <-chan4:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V4, 8)
				case n := <-chan5:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V5, 16)
				case n := <-chan6:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V6, 32)
				case n := <-chan7:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V7, 64)
				case n := <-chan8:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V8, 128)
				case n := <-chan9:
					cont = combineLatestTry9(sink, n, proj, &s, &s.V9, 256)
				}
			}
		})
	}
}

type combineLatestState9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	NBits, CBits uint16

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
	V8 T8
	V9 T9
}

func combineLatestTry9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R, X any](
	sink Observer[R],
	n Notification[X],
	proj func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R,
	s *combineLatestState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	v *X,
	bit uint16,
) bool {
	const FullBits = 511

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			oops := func() { sink.Error(ErrOops) }
			v := Try91(proj, s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7, s.V8, s.V9, oops)
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
