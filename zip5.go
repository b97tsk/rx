package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip5 combines multiple Observables to create an Observable that emits
// projection of values of each of its input Observables.
func Zip5[T1, T2, T3, T4, T5, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
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
		chan3 := make(chan Notification[T3])
		chan4 := make(chan Notification[T4])
		chan5 := make(chan Notification[T5])

		wg.Go(func() {
			var s zipState5[T1, T2, T3, T4, T5]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = zipSink5(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					done = zipSink5(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					done = zipSink5(n, sink, proj, &s, &s.Q3, 4)
				case n := <-chan4:
					done = zipSink5(n, sink, proj, &s, &s.Q4, 8)
				case n := <-chan5:
					done = zipSink5(n, sink, proj, &s, &s.Q5, 16)
				}
			}
		})

		wg.Go(func() { obs1.Subscribe(ctx, chanObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, chanObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, chanObserver(chan3, noop)) })
		wg.Go(func() { obs4.Subscribe(ctx, chanObserver(chan4, noop)) })
		wg.Go(func() { obs5.Subscribe(ctx, chanObserver(chan5, noop)) })
	}
}

type zipState5[T1, T2, T3, T4, T5 any] struct {
	VBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
}

func zipSink5[T1, T2, T3, T4, T5, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5) R,
	s *zipState5[T1, T2, T3, T4, T5],
	q *queue.Queue[X],
	bit uint8,
) bool {
	const FullBits = 31

	switch {
	case n.HasValue:
		q.Push(n.Value)

		if s.VBits |= bit; s.VBits == FullBits {
			var complete bool

			sink.Next(proj(
				zipPop5(s, &s.Q1, 1, &complete),
				zipPop5(s, &s.Q2, 2, &complete),
				zipPop5(s, &s.Q3, 4, &complete),
				zipPop5(s, &s.Q4, 8, &complete),
				zipPop5(s, &s.Q5, 16, &complete),
			))

			if complete {
				sink.Complete()
				return true
			}
		}

	case n.HasError:
		sink.Error(n.Error)
		return true

	default:
		s.CBits |= bit

		if q.Len() == 0 {
			sink.Complete()
			return true
		}
	}

	return false
}

func zipPop5[T1, T2, T3, T4, T5, X any](
	s *zipState5[T1, T2, T3, T4, T5],
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
