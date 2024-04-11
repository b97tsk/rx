package rx

import "github.com/b97tsk/rx/internal/queue"

// ZipWithBuffering3 combines multiple Observables to create an Observable that
// emits mappings of the values emitted by each of its input Observables.
//
// ZipWithBuffering3 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering3[T1, T2, T3, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	mapping func(v1 T1, v2 T2, v3 T3) R,
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
		chan3 := make(chan Notification[T3])

		c.Go(func() {
			var s zipState3[T1, T2, T3]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipTry3(o, n, mapping, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipTry3(o, n, mapping, &s, &s.Q2, 2)
				case n := <-chan3:
					cont = zipTry3(o, n, mapping, &s, &s.Q3, 4)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, ob1, chan1, noop) &&
			subscribeChannel(c, ob2, chan2, noop) &&
			subscribeChannel(c, ob3, chan3, noop)
	}
}

type zipState3[T1, T2, T3 any] struct {
	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
}

func zipTry3[T1, T2, T3, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3) R,
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

			oops := func() { o.Error(ErrOops) }
			v := Try31(
				mapping,
				zipPop3(s, &s.Q1, 1, &complete),
				zipPop3(s, &s.Q2, 2, &complete),
				zipPop3(s, &s.Q3, 4, &complete),
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
