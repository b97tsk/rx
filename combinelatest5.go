package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/waitgroup"
)

// CombineLatest5 combines multiple Observables to create an Observable that
// emits projection of latest values of each of its input Observables.
func CombineLatest5[T1, T2, T3, T4, T5, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
) Observable[R] {
	if proj == nil {
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
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

		ctxHoisted := waitgroup.Hoist(ctx)

		Go(ctxHoisted, func() {
			var s combineLatestState5[T1, T2, T3, T4, T5]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = combineLatestSink5(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = combineLatestSink5(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					done = combineLatestSink5(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					done = combineLatestSink5(n, sink, proj, &s, &s.V4, 8)
				case n := <-chan5:
					done = combineLatestSink5(n, sink, proj, &s, &s.V5, 16)
				}
			}
		})

		Go(ctxHoisted, func() { obs1.Subscribe(ctx, chanObserver(chan1, noop)) })
		Go(ctxHoisted, func() { obs2.Subscribe(ctx, chanObserver(chan2, noop)) })
		Go(ctxHoisted, func() { obs3.Subscribe(ctx, chanObserver(chan3, noop)) })
		Go(ctxHoisted, func() { obs4.Subscribe(ctx, chanObserver(chan4, noop)) })
		Go(ctxHoisted, func() { obs5.Subscribe(ctx, chanObserver(chan5, noop)) })
	}
}

type combineLatestState5[T1, T2, T3, T4, T5 any] struct {
	VBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
}

func combineLatestSink5[T1, T2, T3, T4, T5, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5) R,
	s *combineLatestState5[T1, T2, T3, T4, T5],
	v *X,
	bit uint8,
) bool {
	const FullBits = 31

	switch {
	case n.HasValue:
		*v = n.Value

		if s.VBits |= bit; s.VBits == FullBits {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4, s.V5))
		}

	case n.HasError:
		sink.Error(n.Error)
		return true

	default:
		if s.CBits |= bit; s.CBits == FullBits {
			sink.Complete()
			return true
		}
	}

	return false
}
