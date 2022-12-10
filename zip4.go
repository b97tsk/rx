package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip4 creates an Observable that combines multiple Observables to create
// an Observable that emits projection of values of each of its input
// Observables.
func Zip4[T1, T2, T3, T4, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4) R,
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

		noop := make(chan struct{})

		go func() {
			var s zipState4[T1, T2, T3, T4]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = zipSink4(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					done = zipSink4(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					done = zipSink4(n, sink, proj, &s, &s.Q3, 4)
				case n := <-chan4:
					done = zipSink4(n, sink, proj, &s, &s.Q4, 8)
				}
			}

			close(noop)
		}()

		go obs1.Subscribe(ctx, chanObserver(chan1, noop))
		go obs2.Subscribe(ctx, chanObserver(chan2, noop))
		go obs3.Subscribe(ctx, chanObserver(chan3, noop))
		go obs4.Subscribe(ctx, chanObserver(chan4, noop))
	}
}

type zipState4[T1, T2, T3, T4 any] struct {
	VBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
}

func zipSink4[T1, T2, T3, T4, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4) R,
	s *zipState4[T1, T2, T3, T4],
	q *queue.Queue[X],
	b uint8,
) bool {
	const FullBits = 15

	switch {
	case n.HasValue:
		q.Push(n.Value)

		if s.VBits |= b; s.VBits == FullBits {
			var completed bool

			sink.Next(proj(
				zipPop4(s, &s.Q1, 1, &completed),
				zipPop4(s, &s.Q2, 2, &completed),
				zipPop4(s, &s.Q3, 4, &completed),
				zipPop4(s, &s.Q4, 8, &completed),
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

func zipPop4[T1, T2, T3, T4, X any](
	s *zipState4[T1, T2, T3, T4],
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
