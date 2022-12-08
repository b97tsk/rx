package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip8 creates an Observable that combines multiple Observables to create
// an Observable that emits projection of values of each of its input
// Observables.
func Zip8[T1, T2, T3, T4, T5, T6, T7, T8, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	obs8 Observable[T8],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8) R,
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

		noop := make(chan struct{})

		go func() {
			var s zipState8[T1, T2, T3, T4, T5, T6, T7, T8]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = zipSink8(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					done = zipSink8(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					done = zipSink8(n, sink, proj, &s, &s.Q3, 4)
				case n := <-chan4:
					done = zipSink8(n, sink, proj, &s, &s.Q4, 8)
				case n := <-chan5:
					done = zipSink8(n, sink, proj, &s, &s.Q5, 16)
				case n := <-chan6:
					done = zipSink8(n, sink, proj, &s, &s.Q6, 32)
				case n := <-chan7:
					done = zipSink8(n, sink, proj, &s, &s.Q7, 64)
				case n := <-chan8:
					done = zipSink8(n, sink, proj, &s, &s.Q8, 128)
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
	}
}

type zipState8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	RBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
	Q6 queue.Queue[T6]
	Q7 queue.Queue[T7]
	Q8 queue.Queue[T8]
}

func zipSink8[T1, T2, T3, T4, T5, T6, T7, T8, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7, T8) R,
	s *zipState8[T1, T2, T3, T4, T5, T6, T7, T8],
	q *queue.Queue[X],
	b uint8,
) bool {
	const FullBits = 255

	switch {
	case n.HasValue:
		q.Push(n.Value)

		if s.RBits |= b; s.RBits == FullBits {
			var completed bool

			sink.Next(proj(
				zipPop8(s, &s.Q1, 1, &completed),
				zipPop8(s, &s.Q2, 2, &completed),
				zipPop8(s, &s.Q3, 4, &completed),
				zipPop8(s, &s.Q4, 8, &completed),
				zipPop8(s, &s.Q5, 16, &completed),
				zipPop8(s, &s.Q6, 32, &completed),
				zipPop8(s, &s.Q7, 64, &completed),
				zipPop8(s, &s.Q8, 128, &completed),
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

func zipPop8[T1, T2, T3, T4, T5, T6, T7, T8, X any](
	s *zipState8[T1, T2, T3, T4, T5, T6, T7, T8],
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
