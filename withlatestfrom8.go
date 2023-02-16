package rx

import (
	"context"
)

// WithLatestFrom8 combines the source with 8 other Observables to create an
// Observable that emits projection of the latest values of each Observable,
// only when the source emits.
func WithLatestFrom8[T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	obs8 Observable[T8],
	proj func(v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8) R,
) Operator[T0, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom9(source, obs1, obs2, obs3, obs4, obs5, obs6, obs7, obs8, proj)
		},
	)
}

func withLatestFrom9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](
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

		noop := make(chan struct{})

		go func() {
			var s withLatestFromState9[T1, T2, T3, T4, T5, T6, T7, T8, T9]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V4, 8)
				case n := <-chan5:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V5, 16)
				case n := <-chan6:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V6, 32)
				case n := <-chan7:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V7, 64)
				case n := <-chan8:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V8, 128)
				case n := <-chan9:
					done = withLatestFromSink9(n, sink, proj, &s, &s.V9, 256)
				}
			}

			close(noop)
		}()

		go obs1.Subscribe(ctx, chanObserver(chan1, noop))
		go obs2.Subscribe(ctx, chanObserver(chan2, noop))
		go obs3.Subscribe(ctx, chanObserver(chan3, noop))
		go obs4.Subscribe(ctx, chanObserver(chan4, noop))
		go obs5.Subscribe(ctx, chanObserver(chan5, noop))
		go obs6.Subscribe(ctx, chanObserver(chan6, noop))
		go obs7.Subscribe(ctx, chanObserver(chan7, noop))
		go obs8.Subscribe(ctx, chanObserver(chan8, noop))
		go obs9.Subscribe(ctx, chanObserver(chan9, noop))
	}
}

type withLatestFromState9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	VBits uint16

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

func withLatestFromSink9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R,
	s *withLatestFromState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	v *X,
	b uint16,
) bool {
	const FullBits = 511

	switch {
	case n.HasValue:
		*v = n.Value

		if s.VBits |= b; s.VBits == FullBits && b == 1 {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7, s.V8, s.V9))
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
