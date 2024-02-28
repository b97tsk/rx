package rx

import "github.com/b97tsk/rx/internal/queue"

// ZipWithBuffering2 combines multiple Observables to create an Observable that
// emits projection of values of each of its input Observables.
//
// ZipWithBuffering2 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering2[T1, T2, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	proj func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, sink Observer[R]) {
		c, cancel := c.WithCancel()
		noop := make(chan struct{})
		sink = sink.OnLastNotification(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		c.Go(func() { obs1.Subscribe(c, channelObserver(chan1, noop)) })
		c.Go(func() { obs2.Subscribe(c, channelObserver(chan2, noop)) })

		c.Go(func() {
			var s zipState2[T1, T2]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipTry2(sink, n, proj, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipTry2(sink, n, proj, &s, &s.Q2, 2)
				}
			}
		})
	}
}

type zipState2[T1, T2 any] struct {
	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
}

func zipTry2[T1, T2, R, X any](
	sink Observer[R],
	n Notification[X],
	proj func(T1, T2) R,
	s *zipState2[T1, T2],
	q *queue.Queue[X],
	bit uint8,
) bool {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.NBits |= bit; s.NBits == FullBits {
			var complete bool

			oops := func() { sink.Error(ErrOops) }
			v := Try21(
				proj,
				zipPop2(s, &s.Q1, 1, &complete),
				zipPop2(s, &s.Q2, 2, &complete),
				oops,
			)
			Try1(sink, Next(v), oops)

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

func zipPop2[T1, T2, X any](
	s *zipState2[T1, T2],
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
