package rx

import "github.com/b97tsk/rx/internal/queue"

// ZipWithBuffering2 combines multiple Observables to create an Observable that
// emits mappings of the values emitted by each of its input Observables.
//
// ZipWithBuffering2 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering2[T1, T2, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, cancel := c.WithCancel()
		noop := make(chan struct{})
		o = o.DoOnTermination(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])

		c.Go(func() {
			var s zipState2[T1, T2]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipEmit2(o, n, mapping, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipEmit2(o, n, mapping, &s, &s.Q2, 2)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, ob1, chan1, noop) &&
			subscribeChannel(c, ob2, chan2, noop)
	}
}

type zipState2[T1, T2 any] struct {
	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
}

func zipEmit2[T1, T2, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2) R,
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

			oops := func() { o.Error(ErrOops) }
			v := Try21(
				mapping,
				zipPop2(s, &s.Q1, 1, &complete),
				zipPop2(s, &s.Q2, 2, &complete),
				oops,
			)
			Try1(o, Next(v), oops)

			if complete {
				o.Complete()
				return false
			}
		}

	case KindError:
		o.Error(n.Error)
		return false

	case KindComplete:
		s.CBits |= bit

		if q.Len() == 0 {
			o.Complete()
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
