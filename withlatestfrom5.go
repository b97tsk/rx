package rx

import "context"

// WithLatestFrom5 combines the source with 5 other Observables to create
// an Observable that emits projection of latest values of each Observable,
// only when the source emits.
func WithLatestFrom5[T0, T1, T2, T3, T4, T5, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	proj func(v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
) Operator[T0, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom6(source, obs1, obs2, obs3, obs4, obs5, proj)
		},
	)
}

func withLatestFrom6[T1, T2, T3, T4, T5, T6, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6) R,
) Observable[R] {
	return func(ctx context.Context, sink Observer[R]) {
		wg := WaitGroupFromContext(ctx)
		ctx, cancel := context.WithCancel(ctx)

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

		wg.Go(func() {
			var s withLatestFromState6[T1, T2, T3, T4, T5, T6]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = withLatestFromSink6(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = withLatestFromSink6(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					done = withLatestFromSink6(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					done = withLatestFromSink6(n, sink, proj, &s, &s.V4, 8)
				case n := <-chan5:
					done = withLatestFromSink6(n, sink, proj, &s, &s.V5, 16)
				case n := <-chan6:
					done = withLatestFromSink6(n, sink, proj, &s, &s.V6, 32)
				}
			}
		})

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, channelObserver(chan3, noop)) })
		wg.Go(func() { obs4.Subscribe(ctx, channelObserver(chan4, noop)) })
		wg.Go(func() { obs5.Subscribe(ctx, channelObserver(chan5, noop)) })
		wg.Go(func() { obs6.Subscribe(ctx, channelObserver(chan6, noop)) })
	}
}

type withLatestFromState6[T1, T2, T3, T4, T5, T6 any] struct {
	NBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
}

func withLatestFromSink6[T1, T2, T3, T4, T5, T6, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6) R,
	s *withLatestFromState6[T1, T2, T3, T4, T5, T6],
	v *X,
	bit uint8,
) bool {
	const FullBits = 63

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4, s.V5, s.V6))
		}

	case KindError:
		sink.Error(n.Error)
		return true

	case KindComplete:
		if bit == 1 {
			sink.Complete()
			return true
		}
	}

	return false
}
