package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip9 combines multiple Observables to create an Observable that emits
// projection of values of each of its input Observables.
func Zip9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](
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
		chan6 := make(chan Notification[T6])
		chan7 := make(chan Notification[T7])
		chan8 := make(chan Notification[T8])
		chan9 := make(chan Notification[T9])

		wg.Go(func() {
			var s zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9]

			done := false

			for !done {
				select {
				case n := <-chan1:
					done = zipSink9(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					done = zipSink9(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					done = zipSink9(n, sink, proj, &s, &s.Q3, 4)
				case n := <-chan4:
					done = zipSink9(n, sink, proj, &s, &s.Q4, 8)
				case n := <-chan5:
					done = zipSink9(n, sink, proj, &s, &s.Q5, 16)
				case n := <-chan6:
					done = zipSink9(n, sink, proj, &s, &s.Q6, 32)
				case n := <-chan7:
					done = zipSink9(n, sink, proj, &s, &s.Q7, 64)
				case n := <-chan8:
					done = zipSink9(n, sink, proj, &s, &s.Q8, 128)
				case n := <-chan9:
					done = zipSink9(n, sink, proj, &s, &s.Q9, 256)
				}
			}
		})

		wg.Go(func() { obs1.Subscribe(ctx, chanObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, chanObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, chanObserver(chan3, noop)) })
		wg.Go(func() { obs4.Subscribe(ctx, chanObserver(chan4, noop)) })
		wg.Go(func() { obs5.Subscribe(ctx, chanObserver(chan5, noop)) })
		wg.Go(func() { obs6.Subscribe(ctx, chanObserver(chan6, noop)) })
		wg.Go(func() { obs7.Subscribe(ctx, chanObserver(chan7, noop)) })
		wg.Go(func() { obs8.Subscribe(ctx, chanObserver(chan8, noop)) })
		wg.Go(func() { obs9.Subscribe(ctx, chanObserver(chan9, noop)) })
	}
}

type zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	VBits, CBits uint16

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
	Q6 queue.Queue[T6]
	Q7 queue.Queue[T7]
	Q8 queue.Queue[T8]
	Q9 queue.Queue[T9]
}

func zipSink9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R,
	s *zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	q *queue.Queue[X],
	bit uint16,
) bool {
	const FullBits = 511

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.VBits |= bit; s.VBits == FullBits {
			var complete bool

			sink.Next(proj(
				zipPop9(s, &s.Q1, 1, &complete),
				zipPop9(s, &s.Q2, 2, &complete),
				zipPop9(s, &s.Q3, 4, &complete),
				zipPop9(s, &s.Q4, 8, &complete),
				zipPop9(s, &s.Q5, 16, &complete),
				zipPop9(s, &s.Q6, 32, &complete),
				zipPop9(s, &s.Q7, 64, &complete),
				zipPop9(s, &s.Q8, 128, &complete),
				zipPop9(s, &s.Q9, 256, &complete),
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

func zipPop9[T1, T2, T3, T4, T5, T6, T7, T8, T9, X any](
	s *zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	q *queue.Queue[X],
	bit uint16,
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
