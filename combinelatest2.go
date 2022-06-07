package rx

import (
	"context"
)

// CombineLatest2 creates an Observable that combines multiple Observables to
// create an Observable that emits projection of the latest values of each of
// its input Observables.
func CombineLatest2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v1 T1, v2 T2) R,
) Observable[R] {
	switch {
	case obs1 == nil:
		panic("obs1 == nil")
	case obs2 == nil:
		panic("obs2 == nil")
	case proj == nil:
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		go func() {
			var s combineLatestState2[T1, T2]

			done := ctx.Done()
			exit := false

			for !exit {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case n := <-chan1:
					exit = combineLatestSink2(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					exit = combineLatestSink2(n, sink, proj, &s, &s.V2, 2)
				}
			}
		}()

		go subscribeToChan(ctx, obs1, chan1)
		go subscribeToChan(ctx, obs2, chan2)
	}
}

type combineLatestState2[T1, T2 any] struct {
	RBits, CBits uint8

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

		if s.RBits |= b; s.RBits == FullBits {
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
