package rx

import (
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// ZipWithBuffering8 combines multiple Observables to create an [Observable]
// that emits mappings of the values emitted by each of the input Observables.
//
// ZipWithBuffering8 buffers every value from each input [Observable], which
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
		c, o = Synchronize(c, o)

		var s zipState8[T1, T2, T3, T4, T5, T6, T7, T8]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { zipEmit8(o, n, mapping, &s, &s.q1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { zipEmit8(o, n, mapping, &s, &s.q2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { zipEmit8(o, n, mapping, &s, &s.q3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { zipEmit8(o, n, mapping, &s, &s.q4, 8) }) &&
			ob5.satcc(c, func(n Notification[T5]) { zipEmit8(o, n, mapping, &s, &s.q5, 16) }) &&
			ob6.satcc(c, func(n Notification[T6]) { zipEmit8(o, n, mapping, &s, &s.q6, 32) }) &&
			ob7.satcc(c, func(n Notification[T7]) { zipEmit8(o, n, mapping, &s, &s.q7, 64) }) &&
			ob8.satcc(c, func(n Notification[T8]) { zipEmit8(o, n, mapping, &s, &s.q8, 128) })
	}
}

type zipState8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	mu sync.Mutex

	nbits, cbits uint8

	q1 queue.Queue[T1]
	q2 queue.Queue[T2]
	q3 queue.Queue[T3]
	q4 queue.Queue[T4]
	q5 queue.Queue[T5]
	q6 queue.Queue[T6]
	q7 queue.Queue[T7]
	q8 queue.Queue[T8]
}

func zipEmit8[T1, T2, T3, T4, T5, T6, T7, T8, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5, T6, T7, T8) R,
	s *zipState8[T1, T2, T3, T4, T5, T6, T7, T8],
	q *queue.Queue[X],
	bit uint8,
) {
	const FullBits = 255

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		q.Push(n.Value)

		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits {
			var complete bool

			v := Try81(
				mapping,
				zipPop8(s, &s.q1, 1, &complete),
				zipPop8(s, &s.q2, 2, &complete),
				zipPop8(s, &s.q3, 4, &complete),
				zipPop8(s, &s.q4, 8, &complete),
				zipPop8(s, &s.q5, 16, &complete),
				zipPop8(s, &s.q6, 32, &complete),
				zipPop8(s, &s.q7, 64, &complete),
				zipPop8(s, &s.q8, 128, &complete),
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

func zipPop8[T1, T2, T3, T4, T5, T6, T7, T8, X any](
	s *zipState8[T1, T2, T3, T4, T5, T6, T7, T8],
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
