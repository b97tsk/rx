package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip6 creates an Observable that combines multiple Observables to create
// an Observable that emits projection of values of each of its input
// Observables.
func Zip6[T1, T2, T3, T4, T5, T6, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6) R,
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

		go func() {
			var s zipState6[T1, T2, T3, T4, T5, T6]

			done := ctx.Done()
			exit := false

			for !exit {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case n := <-chan1:
					exit = zipSink6(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					exit = zipSink6(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					exit = zipSink6(n, sink, proj, &s, &s.Q3, 4)
				case n := <-chan4:
					exit = zipSink6(n, sink, proj, &s, &s.Q4, 8)
				case n := <-chan5:
					exit = zipSink6(n, sink, proj, &s, &s.Q5, 16)
				case n := <-chan6:
					exit = zipSink6(n, sink, proj, &s, &s.Q6, 32)
				}
			}
		}()

		go subscribeToChan(ctx, obs1, chan1)
		go subscribeToChan(ctx, obs2, chan2)
		go subscribeToChan(ctx, obs3, chan3)
		go subscribeToChan(ctx, obs4, chan4)
		go subscribeToChan(ctx, obs5, chan5)
		go subscribeToChan(ctx, obs6, chan6)
	}
}

type zipState6[T1, T2, T3, T4, T5, T6 any] struct {
	RBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
	Q6 queue.Queue[T6]
}

func zipSink6[T1, T2, T3, T4, T5, T6, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6) R,
	s *zipState6[T1, T2, T3, T4, T5, T6],
	q *queue.Queue[X],
	b uint8,
) bool {
	const FullBits = 63

	switch {
	case n.HasValue:
		q.Push(n.Value)

		if s.RBits |= b; s.RBits == FullBits {
			var completed bool

			sink.Next(proj(
				zipPop6(s, &s.Q1, 1, &completed),
				zipPop6(s, &s.Q2, 2, &completed),
				zipPop6(s, &s.Q3, 4, &completed),
				zipPop6(s, &s.Q4, 8, &completed),
				zipPop6(s, &s.Q5, 16, &completed),
				zipPop6(s, &s.Q6, 32, &completed),
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

func zipPop6[T1, T2, T3, T4, T5, T6, X any](
	s *zipState6[T1, T2, T3, T4, T5, T6],
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
