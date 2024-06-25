package rx

import (
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering5 combines multiple Observables to create an Observable that
// emits mappings of the values emitted by each of its input Observables.
//
// ZipWithBuffering5 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering5[T1, T2, T3, T4, T5, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s zipState5[T1, T2, T3, T4, T5]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { zipEmit5(o, n, mapping, &s, &s.Q1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { zipEmit5(o, n, mapping, &s, &s.Q2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { zipEmit5(o, n, mapping, &s, &s.Q3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { zipEmit5(o, n, mapping, &s, &s.Q4, 8) }) &&
			ob5.satcc(c, func(n Notification[T5]) { zipEmit5(o, n, mapping, &s, &s.Q5, 16) })
	}
}

type zipState5[T1, T2, T3, T4, T5 any] struct {
	sync.Mutex

	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
	Q4 queue.Queue[T4]
	Q5 queue.Queue[T5]
}

func zipEmit5[T1, T2, T3, T4, T5, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5) R,
	s *zipState5[T1, T2, T3, T4, T5],
	q *queue.Queue[X],
	bit uint8,
) {
	const FullBits = 31

	switch n.Kind {
	case KindNext:
		s.Lock()
		q.Push(n.Value)

		nbits := s.NBits
		nbits |= bit
		s.NBits = nbits

		if nbits == FullBits {
			var complete bool

			v := Try51(
				mapping,
				zipPop5(s, &s.Q1, 1, &complete),
				zipPop5(s, &s.Q2, 2, &complete),
				zipPop5(s, &s.Q3, 4, &complete),
				zipPop5(s, &s.Q4, 8, &complete),
				zipPop5(s, &s.Q5, 16, &complete),
				s.Unlock,
			)
			s.Unlock()
			o.Next(v)

			if complete {
				o.Complete()
			}

			return
		}

		s.Unlock()

	case KindError:
		o.Error(n.Error)

	case KindComplete:
		s.Lock()
		s.CBits |= bit
		complete := q.Len() == 0
		s.Unlock()

		if complete {
			o.Complete()
		}
	}
}

func zipPop5[T1, T2, T3, T4, T5, X any](
	s *zipState5[T1, T2, T3, T4, T5],
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
