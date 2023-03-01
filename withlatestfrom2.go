package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/waitgroup"
)

// WithLatestFrom2 combines the source with 2 other Observables to create
// an Observable that emits projection of latest values of each Observable,
// only when the source emits.
func WithLatestFrom2[T0, T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v0 T0, v1 T1, v2 T2) R,
) Operator[T0, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom3(source, obs1, obs2, proj)
		},
	)
}

func withLatestFrom3[T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v1 T1, v2 T2, v3 T3) R,
) Observable[R] {
	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		noop := make(chan struct{})

		sink = sink.WithCancel(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])

		ctxHoisted := waitgroup.Hoist(ctx)

		Go(ctxHoisted, func() {
			var s withLatestFromState3[T1, T2, T3]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = withLatestFromSink3(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = withLatestFromSink3(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					done = withLatestFromSink3(n, sink, proj, &s, &s.V3, 4)
				}
			}
		})

		Go(ctxHoisted, func() { obs1.Subscribe(ctx, chanObserver(chan1, noop)) })
		Go(ctxHoisted, func() { obs2.Subscribe(ctx, chanObserver(chan2, noop)) })
		Go(ctxHoisted, func() { obs3.Subscribe(ctx, chanObserver(chan3, noop)) })
	}
}

type withLatestFromState3[T1, T2, T3 any] struct {
	VBits uint8

	V1 T1
	V2 T2
	V3 T3
}

func withLatestFromSink3[T1, T2, T3, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3) R,
	s *withLatestFromState3[T1, T2, T3],
	v *X,
	b uint8,
) bool {
	const FullBits = 7

	switch {
	case n.HasValue:
		*v = n.Value

		if s.VBits |= b; s.VBits == FullBits && b == 1 {
			sink.Next(proj(s.V1, s.V2, s.V3))
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
