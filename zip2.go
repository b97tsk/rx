package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// Zip2 creates an Observable that combines multiple Observables to create
// an Observable that emits projection of values of each of its input
// Observables.
func Zip2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v1 T1, v2 T2) R,
) Observable[R] {
	if proj == nil {
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
		ctx, cancel := context.WithCancel(ctx)

		noop := make(chan struct{})

		sink = sink.WithCancel(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		ctxHoisted := waitgroup.Hoist(ctx)

		Go(ctxHoisted, func() {
			var s zipState2[T1, T2]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = zipSink2(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					done = zipSink2(n, sink, proj, &s, &s.Q2, 2)
				}
			}
		})

		Go(ctxHoisted, func() { obs1.Subscribe(ctx, chanObserver(chan1, noop)) })
		Go(ctxHoisted, func() { obs2.Subscribe(ctx, chanObserver(chan2, noop)) })
	}
}

type zipState2[T1, T2 any] struct {
	VBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
}

func zipSink2[T1, T2, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2) R,
	s *zipState2[T1, T2],
	q *queue.Queue[X],
	b uint8,
) bool {
	const FullBits = 3

	switch {
	case n.HasValue:
		q.Push(n.Value)

		if s.VBits |= b; s.VBits == FullBits {
			var completed bool

			sink.Next(proj(
				zipPop2(s, &s.Q1, 1, &completed),
				zipPop2(s, &s.Q2, 2, &completed),
			))

			if completed {
				sink.Complete()
				return true
			}
		}

	case n.HasError:
		sink.Error(n.Error)
		return true

	default:
		s.CBits |= b

		if q.Len() == 0 {
			sink.Complete()
			return true
		}
	}

	return false
}

func zipPop2[T1, T2, X any](
	s *zipState2[T1, T2],
	q *queue.Queue[X],
	b uint8,
	completed *bool,
) X {
	v := q.Pop()

	if q.Len() == 0 {
		s.VBits &^= b

		if s.CBits&b != 0 {
			*completed = true
		}
	}

	return v
}
