package rx

import "github.com/b97tsk/rx/internal/queue"

// ZipWithBuffering8 combines multiple Observables to create an Observable that
// emits mappings of the values emitted by each of its input Observables.
//
// ZipWithBuffering8 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering8[T1, T2, T3, T4, T5, T6, T7, T8, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	ob6 Observable[T6],
	ob7 Observable[T7],
	ob8 Observable[T8],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8) R,
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
		chan5 := make(chan Notification[T5])
		chan6 := make(chan Notification[T6])
		chan7 := make(chan Notification[T7])
		chan8 := make(chan Notification[T8])

		c.Go(func() {
			var s zipState8[T1, T2, T3, T4, T5, T6, T7, T8]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = zipEmit8(o, n, mapping, &s, &s.Q1, 1)
				case n := <-chan2:
					cont = zipEmit8(o, n, mapping, &s, &s.Q2, 2)
				case n := <-chan3:
					cont = zipEmit8(o, n, mapping, &s, &s.Q3, 4)
				case n := <-chan4:
					cont = zipEmit8(o, n, mapping, &s, &s.Q4, 8)
				case n := <-chan5:
					cont = zipEmit8(o, n, mapping, &s, &s.Q5, 16)
				case n := <-chan6:
					cont = zipEmit8(o, n, mapping, &s, &s.Q6, 32)
				case n := <-chan7:
					cont = zipEmit8(o, n, mapping, &s, &s.Q7, 64)
				case n := <-chan8:
					cont = zipEmit8(o, n, mapping, &s, &s.Q8, 128)
				}
			}
		})

		_ = true &&
			subscribeChannel(c, ob1, chan1, noop) &&
			subscribeChannel(c, ob2, chan2, noop) &&
			subscribeChannel(c, ob3, chan3, noop) &&
			subscribeChannel(c, ob4, chan4, noop) &&
			subscribeChannel(c, ob5, chan5, noop) &&
			subscribeChannel(c, ob6, chan6, noop) &&
			subscribeChannel(c, ob7, chan7, noop) &&
			subscribeChannel(c, ob8, chan8, noop)
	}
}

type zipState8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
	Q6 queue.Queue[T6]
	Q7 queue.Queue[T7]
	Q8 queue.Queue[T8]
}

func zipEmit8[T1, T2, T3, T4, T5, T6, T7, T8, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5, T6, T7, T8) R,
	s *zipState8[T1, T2, T3, T4, T5, T6, T7, T8],
	q *queue.Queue[X],
	bit uint8,
) bool {
	const FullBits = 255

	switch n.Kind {
	case KindNext:
		q.Push(n.Value)

		if s.NBits |= bit; s.NBits == FullBits {
			var complete bool

			oops := func() { o.Error(ErrOops) }
			v := Try81(
				mapping,
				zipPop8(s, &s.Q1, 1, &complete),
				zipPop8(s, &s.Q2, 2, &complete),
				zipPop8(s, &s.Q3, 4, &complete),
				zipPop8(s, &s.Q4, 8, &complete),
				zipPop8(s, &s.Q5, 16, &complete),
				zipPop8(s, &s.Q6, 32, &complete),
				zipPop8(s, &s.Q7, 64, &complete),
				zipPop8(s, &s.Q8, 128, &complete),
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

func zipPop8[T1, T2, T3, T4, T5, T6, T7, T8, X any](
	s *zipState8[T1, T2, T3, T4, T5, T6, T7, T8],
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
