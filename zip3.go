package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip3 creates an Observable that combines multiple Observables to create
// an Observable that emits projection of values of each of its input
// Observables.
func Zip3[T1, T2, T3, R any](
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
			var s zipState3[T1, T2, T3]

			done := ctx.Done()
			exit := false

			for !exit {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case n := <-chan1:
					exit = zipSink3(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					exit = zipSink3(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					exit = zipSink3(n, sink, proj, &s, &s.Q3, 4)
				}
			}
		}()

		go subscribeToChan(ctx, obs1, chan1)
		go subscribeToChan(ctx, obs2, chan2)
		go subscribeToChan(ctx, obs3, chan3)
	}
}

type zipState3[T1, T2, T3 any] struct {
	RBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
}

func zipSink3[T1, T2, T3, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3) R,
	s *zipState3[T1, T2, T3],
	q *queue.Queue[X],
	b uint8,
) bool {
	const FullBits = 7

	switch {
	case n.HasValue:
		q.Push(n.Value)

		if s.RBits |= b; s.RBits == FullBits {
			var completed bool

			sink.Next(proj(
				zipPop3(s, &s.Q1, 1, &completed),
				zipPop3(s, &s.Q2, 2, &completed),
				zipPop3(s, &s.Q3, 4, &completed),
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

func zipPop3[T1, T2, T3, X any](
	s *zipState3[T1, T2, T3],
	q *queue.Queue[X],
	b uint8,
	completed *bool,
) X {
	v := q.Pop()

	if q.Len() == 0 {
		s.RBits &^= b

		if s.CBits&b != 0 {
			*completed = true
		}
	}

	return v
}
