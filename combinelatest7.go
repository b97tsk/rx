package rx

import (
	"context"
)

// CombineLatest7 creates an Observable that combines multiple Observables to
// create an Observable that emits projection of the latest values of each of
// its input Observables.
func CombineLatest7[T1, T2, T3, T4, T5, T6, T7, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
) Observable[R] {
	switch {
	case obs1 == nil:
		panic("obs1 == nil")
	case obs2 == nil:
		panic("obs2 == nil")
	case obs3 == nil:
		panic("obs3 == nil")
	case obs4 == nil:
		panic("obs4 == nil")
	case obs5 == nil:
		panic("obs5 == nil")
	case obs6 == nil:
		panic("obs6 == nil")
	case obs7 == nil:
		panic("obs7 == nil")
	case proj == nil:
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])
		chan4 := make(chan Notification[T4])
		chan5 := make(chan Notification[T5])
		chan6 := make(chan Notification[T6])
		chan7 := make(chan Notification[T7])

		go func() {
			var s combineLatestState7[T1, T2, T3, T4, T5, T6, T7]

			done := ctx.Done()
			exit := false

			for !exit {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case n := <-chan1:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V4, 8)
				case n := <-chan5:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V5, 16)
				case n := <-chan6:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V6, 32)
				case n := <-chan7:
					exit = combineLatestSink7(n, sink, proj, &s, &s.V7, 64)
				}
			}
		}()

		go subscribeToChan(ctx, obs1, chan1)
		go subscribeToChan(ctx, obs2, chan2)
		go subscribeToChan(ctx, obs3, chan3)
		go subscribeToChan(ctx, obs4, chan4)
		go subscribeToChan(ctx, obs5, chan5)
		go subscribeToChan(ctx, obs6, chan6)
		go subscribeToChan(ctx, obs7, chan7)
	}
}

type combineLatestState7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	RBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
}

func combineLatestSink7[T1, T2, T3, T4, T5, T6, T7, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7) R,
	s *combineLatestState7[T1, T2, T3, T4, T5, T6, T7],
	v *X,
	b uint8,
) bool {
	const FullBits = 127

	switch {
	case n.HasValue:
		*v = n.Value

		if s.RBits |= b; s.RBits == FullBits {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7))
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
