package rx

import (
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering9 combines multiple Observables to create an Observable that
// emits mappings of the values emitted by each of its input Observables.
//
// ZipWithBuffering9 buffers every value from each input Observable, which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	ob6 Observable[T6],
	ob7 Observable[T7],
	ob8 Observable[T8],
	ob9 Observable[T9],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { zipEmit9(o, n, mapping, &s, &s.Q1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { zipEmit9(o, n, mapping, &s, &s.Q2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { zipEmit9(o, n, mapping, &s, &s.Q3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { zipEmit9(o, n, mapping, &s, &s.Q4, 8) }) &&
			ob5.satcc(c, func(n Notification[T5]) { zipEmit9(o, n, mapping, &s, &s.Q5, 16) }) &&
			ob6.satcc(c, func(n Notification[T6]) { zipEmit9(o, n, mapping, &s, &s.Q6, 32) }) &&
			ob7.satcc(c, func(n Notification[T7]) { zipEmit9(o, n, mapping, &s, &s.Q7, 64) }) &&
			ob8.satcc(c, func(n Notification[T8]) { zipEmit9(o, n, mapping, &s, &s.Q8, 128) }) &&
			ob9.satcc(c, func(n Notification[T9]) { zipEmit9(o, n, mapping, &s, &s.Q9, 256) })
	}
}

type zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	sync.Mutex

	NBits, CBits uint16

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

func zipEmit9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5, T6, T7, T8, T9) R,
	s *zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	q *queue.Queue[X],
	bit uint16,
) {
	const FullBits = 511

	switch n.Kind {
	case KindNext:
		s.Lock()
		q.Push(n.Value)

		nbits := s.NBits
		nbits |= bit
		s.NBits = nbits

		if nbits == FullBits {
			var complete bool

			v := Try91(
				mapping,
				zipPop9(s, &s.Q1, 1, &complete),
				zipPop9(s, &s.Q2, 2, &complete),
				zipPop9(s, &s.Q3, 4, &complete),
				zipPop9(s, &s.Q4, 8, &complete),
				zipPop9(s, &s.Q5, 16, &complete),
				zipPop9(s, &s.Q6, 32, &complete),
				zipPop9(s, &s.Q7, 64, &complete),
				zipPop9(s, &s.Q8, 128, &complete),
				zipPop9(s, &s.Q9, 256, &complete),
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

func zipPop9[T1, T2, T3, T4, T5, T6, T7, T8, T9, X any](
	s *zipState9[T1, T2, T3, T4, T5, T6, T7, T8, T9],
	q *queue.Queue[X],
	bit uint16,
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
