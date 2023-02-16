package rx

import (
	"context"
)

// WithLatestFrom1 combines the source with another Observable to create an
// Observable that emits projection of the latest values of each Observable,
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
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		noop := make(chan struct{})

		go func() {
			var s withLatestFromState2[T1, T2]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = withLatestFromSink2(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					done = withLatestFromSink2(n, sink, proj, &s, &s.V2, 2)
				}
			}

			close(noop)
		}()

		go obs1.Subscribe(ctx, chanObserver(chan1, noop))
		go obs2.Subscribe(ctx, chanObserver(chan2, noop))
	}
}

type withLatestFromState2[T1, T2 any] struct {
	VBits uint8

	V1 T1
	V2 T2
}

func withLatestFromSink2[T1, T2, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2) R,
	s *withLatestFromState2[T1, T2],
	v *X,
	b uint8,
) bool {
	const FullBits = 3

	switch {
	case n.HasValue:
		*v = n.Value

		if s.VBits |= b; s.VBits == FullBits && b == 1 {
			sink.Next(proj(s.V1, s.V2))
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
