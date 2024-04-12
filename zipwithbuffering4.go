package rx

import "github.com/b97tsk/rx/internal/queue"

// ZipWithBuffering4 combines multiple Observables to create an Observable that
// emits mappings of the values emitted by each of its input Observables.
//
// ZipWithBuffering4 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering4[T1, T2, T3, T4, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4) R,
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
		chan4 := make(chan Notification[T4])

		c.Go(func() {
			var s zipState4[T1, T2, T3, T4]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipEmit4(o, n, mapping, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipEmit4(o, n, mapping, &s, &s.Q2, 2)
				case n := <-chan3:
					cont = zipEmit4(o, n, mapping, &s, &s.Q3, 4)
				case n := <-chan4:
					cont = zipEmit4(o, n, mapping, &s, &s.Q4, 8)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, ob1, chan1, noop) &&
			subscribeChannel(c, ob2, chan2, noop) &&
			subscribeChannel(c, ob3, chan3, noop) &&
			subscribeChannel(c, ob4, chan4, noop)
	}
}

type zipState4[T1, T2, T3, T4 any] struct {
	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
}

func zipEmit4[T1, T2, T3, T4, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4) R,
	s *zipState4[T1, T2, T3, T4],
	q *queue.Queue[X],
	bit uint8,
) bool {
	const FullBits = 15

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.NBits |= bit; s.NBits == FullBits {
			var complete bool

			oops := func() { o.Error(ErrOops) }
			v := Try41(
				mapping,
				zipPop4(s, &s.Q1, 1, &complete),
				zipPop4(s, &s.Q2, 2, &complete),
				zipPop4(s, &s.Q3, 4, &complete),
				zipPop4(s, &s.Q4, 8, &complete),
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

func zipPop4[T1, T2, T3, T4, X any](
	s *zipState4[T1, T2, T3, T4],
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
