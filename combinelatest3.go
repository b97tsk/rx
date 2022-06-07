package rx

import (
	"context"
)

// CombineLatest3 creates an Observable that combines multiple Observables to
// create an Observable that emits projection of the latest values of each of
// its input Observables.
func CombineLatest3[T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v1 T1, v2 T2, v3 T3) R,
) Observable[R] {
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

	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])

		go func() {
			var s combineLatestState3[T1, T2, T3]

			done := ctx.Done()
			exit := false

			for !exit {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case n := <-chan1:
					exit = combineLatestSink3(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					exit = combineLatestSink3(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					exit = combineLatestSink3(n, sink, proj, &s, &s.V3, 4)
				}
			}
		}()

		go subscribeToChan(ctx, obs1, chan1)
		go subscribeToChan(ctx, obs2, chan2)
		go subscribeToChan(ctx, obs3, chan3)
	}
}

type combineLatestState3[T1, T2, T3 any] struct {
	RBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
}

func combineLatestSink3[T1, T2, T3, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3) R,
	s *combineLatestState3[T1, T2, T3],
	v *X,
	b uint8,
) bool {
	const FullBits = 7

	switch {
	case n.HasValue:
		*v = n.Value

		if s.RBits |= b; s.RBits == FullBits {
			sink.Next(proj(s.V1, s.V2, s.V3))
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
