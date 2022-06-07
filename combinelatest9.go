package rx

import (
	"context"
)

// CombineLatest9 creates an Observable that combines multiple Observables to
// create an Observable that emits projection of the latest values of each of
// its input Observables.
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
	case obs8 == nil:
		panic("obs8 == nil")
	case obs9 == nil:
		panic("obs9 == nil")
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
		chan8 := make(chan Notification[T8])
		chan9 := make(chan Notification[T9])

		go func() {
			var s combineLatestState9[T1, T2, T3, T4, T5, T6, T7, T8, T9]

			done := ctx.Done()
			exit := false

			for !exit {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case n := <-chan1:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V4, 8)
				case n := <-chan5:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V5, 16)
				case n := <-chan6:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V6, 32)
				case n := <-chan7:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V7, 64)
				case n := <-chan8:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V8, 128)
				case n := <-chan9:
					exit = combineLatestSink9(n, sink, proj, &s, &s.V9, 256)
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
		go subscribeToChan(ctx, obs8, chan8)
		go subscribeToChan(ctx, obs9, chan9)
	}
}

type combineLatestState9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	RBits, CBits uint16

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

func combineLatestSink9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R,
	s *combineLatestState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	v *X,
	b uint16,
) bool {
	const FullBits = 511

	switch {
	case n.HasValue:
		*v = n.Value

		if s.RBits |= b; s.RBits == FullBits {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7, s.V8, s.V9))
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
