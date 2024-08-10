package rx

import (
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering2 combines multiple Observables to create an [Observable]
// that emits mappings of the values emitted by each of the input Observables.
//
// ZipWithBuffering2 buffers every value from each input [Observable], which
// might consume a lot of memory over time if there are lots of values emitting
// faster than zipping.
func ZipWithBuffering2[T1, T2, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Synchronize(c, o)

		var s zipState2[T1, T2]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { zipEmit2(o, n, mapping, &s, &s.q1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { zipEmit2(o, n, mapping, &s, &s.q2, 2) })
	}
}

type zipState2[T1, T2 any] struct {
	mu sync.Mutex

	nbits, cbits uint8

	q1 queue.Queue[T1]
	q2 queue.Queue[T2]
}

func zipEmit2[T1, T2, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2) R,
	s *zipState2[T1, T2],
	q *queue.Queue[X],
	bit uint8,
) {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		q.Push(n.Value)

		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits {
			var complete bool

			v := Try21(
				mapping,
				zipPop2(s, &s.q1, 1, &complete),
				zipPop2(s, &s.q2, 2, &complete),
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

func zipPop2[T1, T2, X any](
	s *zipState2[T1, T2],
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
