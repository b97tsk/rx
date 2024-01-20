package rx

import "context"

// WithLatestFrom1 combines the source with another Observable to create
// an Observable that emits projection of latest values of each Observable,
// only when the source emits.
func WithLatestFrom1[T0, T1, R any](
	obs1 Observable[T1],
	proj func(v0 T0, v1 T1) R,
) Operator[T0, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom2(source, obs1, proj)
		},
	)
}

func withLatestFrom2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v1 T1, v2 T2) R,
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

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })

		wg.Go(func() {
			var s withLatestFromState2[T1, T2]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = withLatestFromSink2(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					cont = withLatestFromSink2(n, sink, proj, &s, &s.V2, 2)
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

func withLatestFromSink2[T1, T2, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2) R,
	s *withLatestFromState2[T1, T2],
	v *X,
	bit uint8,
) bool {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits && bit == 1 {
			sink.Next(proj(s.V1, s.V2))
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
