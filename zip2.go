package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip2 combines multiple Observables to create an Observable that emits
// projection of values of each of its input Observables.
func Zip2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v1 T1, v2 T2) R,
) Observable[R] {
	if proj == nil {
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
		wg := WaitGroupFromContext(ctx)
		ctx, cancel := context.WithCancel(ctx)

		noop := make(chan struct{})

		sink = sink.OnLastNotification(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		wg.Go(func() {
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

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })
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
	bit uint8,
) bool {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.VBits |= bit; s.VBits == FullBits {
			var complete bool

			sink.Next(proj(
				zipPop2(s, &s.Q1, 1, &complete),
				zipPop2(s, &s.Q2, 2, &complete),
			))

			if complete {
				sink.Complete()
				return true
			}
		}

	case KindError:
		sink.Error(n.Error)
		return true

	case KindComplete:
		s.CBits |= bit

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
	bit uint8,
	complete *bool,
) X {
	v := q.Pop()

	if q.Len() == 0 {
		s.VBits &^= bit

		if s.CBits&bit != 0 {
			*complete = true
		}
	}

	return v
}
