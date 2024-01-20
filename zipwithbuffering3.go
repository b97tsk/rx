package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering3 combines multiple Observables to create an Observable that
// emits projection of values of each of its input Observables.
//
// ZipWithBuffering3 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering3[T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v1 T1, v2 T2, v3 T3) R,
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

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, channelObserver(chan3, noop)) })

		wg.Go(func() {
			var s zipState3[T1, T2, T3]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipSink3(n, sink, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipSink3(n, sink, proj, &s, &s.Q2, 2)
				case n := <-chan3:
					cont = zipSink3(n, sink, proj, &s, &s.Q3, 4)
				}
			}
		})
	}
}

type zipState3[T1, T2, T3 any] struct {
	NBits, CBits uint8

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
	bit uint8,
) bool {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.NBits |= bit; s.NBits == FullBits {
			var complete bool

			sink.Next(proj(
				zipPop3(s, &s.Q1, 1, &complete),
				zipPop3(s, &s.Q2, 2, &complete),
				zipPop3(s, &s.Q3, 4, &complete),
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

func zipPop3[T1, T2, T3, X any](
	s *zipState3[T1, T2, T3],
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
