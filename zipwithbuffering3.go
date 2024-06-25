package rx

import (
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

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
		c, o = Serialize(c, o)

		var s zipState3[T1, T2, T3]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { zipEmit3(o, n, mapping, &s, &s.Q1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { zipEmit3(o, n, mapping, &s, &s.Q2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { zipEmit3(o, n, mapping, &s, &s.Q3, 4) })
	}
}

type zipState3[T1, T2, T3 any] struct {
	sync.Mutex

	NBits, CBits uint8

	Q1 queue.Queue[T1]
	Q2 queue.Queue[T2]
	Q3 queue.Queue[T3]
}

func zipEmit3[T1, T2, T3, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3) R,
	s *zipState3[T1, T2, T3],
	q *queue.Queue[X],
	bit uint8,
) {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		s.Lock()
		q.Push(n.Value)

		nbits := s.NBits
		nbits |= bit
		s.NBits = nbits

		if nbits == FullBits {
			var complete bool

			v := Try31(
				mapping,
				zipPop3(s, &s.Q1, 1, &complete),
				zipPop3(s, &s.Q2, 2, &complete),
				zipPop3(s, &s.Q3, 4, &complete),
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
