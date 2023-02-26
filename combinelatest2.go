package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/waitgroup"
)

// CombineLatest2 creates an Observable that combines multiple Observables to
// create an Observable that emits projection of the latest values of each of
// its input Observables.
func CombineLatest2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v1 T1, v2 T2) R,
) Observable[R] {
	if proj == nil {
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		noop := make(chan struct{})

		ctxHoisted := waitgroup.Hoist(ctx)

		Go(ctxHoisted, func() {
			var s combineLatestState2[T1, T2]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = combineLatestSink2(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = combineLatestSink2(n, sink, proj, &s, &s.V2, 2)
				}
			}

			close(noop)
		})

		Go(ctxHoisted, func() { obs1.Subscribe(ctx, chanObserver(chan1, noop)) })
		Go(ctxHoisted, func() { obs2.Subscribe(ctx, chanObserver(chan2, noop)) })
	}
}

type combineLatestState2[T1, T2 any] struct {
	VBits, CBits uint8

	V1 T1
	V2 T2
}

func combineLatestSink2[T1, T2, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2) R,
	s *combineLatestState2[T1, T2],
	v *X,
	b uint8,
) bool {
	const FullBits = 3

	switch {
	case n.HasValue:
		*v = n.Value

		if s.VBits |= b; s.VBits == FullBits {
			sink.Next(proj(s.V1, s.V2))
		}

	case n.HasError:
		sink.Error(n.Error)
		return true

	default:
		if s.CBits |= b; s.CBits == FullBits {
			sink.Complete()
			return true
		}
	}

	return false
}
