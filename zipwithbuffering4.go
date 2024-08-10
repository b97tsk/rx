package rx

import (
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering4 combines multiple Observables to create an [Observable]
// that emits mappings of the values emitted by each of the input Observables.
//
// ZipWithBuffering4 buffers every value from each input [Observable], which
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
		c, o = Synchronize(c, o)

		var s zipState4[T1, T2, T3, T4]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { zipEmit4(o, n, mapping, &s, &s.q1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { zipEmit4(o, n, mapping, &s, &s.q2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { zipEmit4(o, n, mapping, &s, &s.q3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { zipEmit4(o, n, mapping, &s, &s.q4, 8) })
	}
}

type zipState4[T1, T2, T3, T4 any] struct {
	mu sync.Mutex

	nbits, cbits uint8

	q1 queue.Queue[T1]
	q2 queue.Queue[T2]
	q3 queue.Queue[T3]
	q4 queue.Queue[T4]
}

func zipEmit4[T1, T2, T3, T4, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4) R,
	s *zipState4[T1, T2, T3, T4],
	q *queue.Queue[X],
	bit uint8,
) {
	const FullBits = 15

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		q.Push(n.Value)

		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits {
			var complete bool

			v := Try41(
				mapping,
				zipPop4(s, &s.q1, 1, &complete),
				zipPop4(s, &s.q2, 2, &complete),
				zipPop4(s, &s.q3, 4, &complete),
				zipPop4(s, &s.q4, 8, &complete),
				s.mu.Unlock,
			)
			s.mu.Unlock()
			o.Next(v)

			if complete {
				o.Complete()
			}

			return
		}

		s.mu.Unlock()

	case KindComplete:
		s.mu.Lock()
		s.cbits |= bit
		complete := q.Len() == 0
		s.mu.Unlock()

		if complete {
			o.Complete()
		}

	case KindError:
		o.Error(n.Error)

	case KindStop:
		o.Stop(n.Error)
	}
}

func zipPop4[T1, T2, T3, T4, X any](
	s *zipState4[T1, T2, T3, T4],
	q *queue.Queue[X],
	bit uint8,
	complete *bool,
) X {
	v := q.Pop()

	if q.Len() == 0 {
		s.nbits &^= bit

		if s.cbits&bit != 0 {
			*complete = true
		}
	}

	return v
}
