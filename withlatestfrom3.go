package rx

import (
	"context"
)

// WithLatestFrom3 combines the source with 3 other Observables to create an
// Observable that emits projection of the latest values of each Observable,
// only when the source emits.
func WithLatestFrom3[T0, T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v0 T0, v1 T1, v2 T2, v3 T3) R,
) Operator[T0, R] {
	switch {
	case obs1 == nil:
		panic("obs1 == nil")
	case obs2 == nil:
		panic("obs2 == nil")
	case obs3 == nil:
		panic("obs3 == nil")
	case proj == nil:
		panic("proj == nil")
	}

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
	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])
		chan4 := make(chan Notification[T4])

		noop := make(chan struct{})

		go func() {
			var s withLatestFromState4[T1, T2, T3, T4]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = withLatestFromSink4(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = withLatestFromSink4(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					done = withLatestFromSink4(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					done = withLatestFromSink4(n, sink, proj, &s, &s.V4, 8)
				}
			}

			close(noop)
		}()

		go obs1.Subscribe(ctx, chanObserver(chan1, noop))
		go obs2.Subscribe(ctx, chanObserver(chan2, noop))
		go obs3.Subscribe(ctx, chanObserver(chan3, noop))
		go obs4.Subscribe(ctx, chanObserver(chan4, noop))
	}
}

type withLatestFromState4[T1, T2, T3, T4 any] struct {
	VBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

func withLatestFromSink4[T1, T2, T3, T4, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4) R,
	s *withLatestFromState4[T1, T2, T3, T4],
	v *X,
	b uint8,
) bool {
	const FullBits = 15

	switch {
	case n.HasValue:
		*v = n.Value

		if s.VBits |= b; s.VBits == FullBits && b == 1 {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4))
		}

	case n.HasError:
		sink.Error(n.Error)
		return true

	default:
		if b == 1 {
			sink.Complete()
			return true
		}
	}

	return false
}
