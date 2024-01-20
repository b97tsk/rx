package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering7 combines multiple Observables to create an Observable that
// emits projection of values of each of its input Observables.
//
// ZipWithBuffering7 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering7[T1, T2, T3, T4, T5, T6, T7, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
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

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, channelObserver(chan3, noop)) })
		wg.Go(func() { obs4.Subscribe(ctx, channelObserver(chan4, noop)) })
		wg.Go(func() { obs5.Subscribe(ctx, channelObserver(chan5, noop)) })
		wg.Go(func() { obs6.Subscribe(ctx, channelObserver(chan6, noop)) })
		wg.Go(func() { obs7.Subscribe(ctx, channelObserver(chan7, noop)) })

		wg.Go(func() {
			var s zipState7[T1, T2, T3, T4, T5, T6, T7]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipSink7(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipSink7(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					cont = zipSink7(n, sink, proj, &s, &s.Q3, 4)
				case n := <-chan4:
					cont = zipSink7(n, sink, proj, &s, &s.Q4, 8)
				case n := <-chan5:
					cont = zipSink7(n, sink, proj, &s, &s.Q5, 16)
				case n := <-chan6:
					cont = zipSink7(n, sink, proj, &s, &s.Q6, 32)
				case n := <-chan7:
					cont = zipSink7(n, sink, proj, &s, &s.Q7, 64)
				}
			}
		})
	}
}

type zipState7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
	Q6 queue.Queue[T6]
	Q7 queue.Queue[T7]
}

func zipSink7[T1, T2, T3, T4, T5, T6, T7, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7) R,
	s *zipState7[T1, T2, T3, T4, T5, T6, T7],
	q *queue.Queue[X],
	bit uint8,
) bool {
	const FullBits = 127

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.NBits |= bit; s.NBits == FullBits {
			var complete bool

			sink.Next(proj(
				zipPop7(s, &s.Q1, 1, &complete),
				zipPop7(s, &s.Q2, 2, &complete),
				zipPop7(s, &s.Q3, 4, &complete),
				zipPop7(s, &s.Q4, 8, &complete),
				zipPop7(s, &s.Q5, 16, &complete),
				zipPop7(s, &s.Q6, 32, &complete),
				zipPop7(s, &s.Q7, 64, &complete),
			))

			if complete {
				sink.Complete()
				return false
			}
		}

	case KindError:
		sink.Error(n.Error)
		return false

	case KindComplete:
		s.CBits |= bit

		if q.Len() == 0 {
			sink.Complete()
			return false
		}
	}

	return true
}

func zipPop7[T1, T2, T3, T4, T5, T6, T7, X any](
	s *zipState7[T1, T2, T3, T4, T5, T6, T7],
	q *queue.Queue[X],
	bit uint8,
	complete *bool,
) X {
	v := q.Pop()

	if q.Len() == 0 {
		s.NBits &^= bit

		if s.CBits&bit != 0 {
			*complete = true
		}
	}

	return v
}
