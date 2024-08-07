package rx

import "sync"

// WithLatestFrom1 combines the source with another Observable to create
// an Observable that emits mappings of the latest values emitted by each
// Observable, only when the source emits.
func WithLatestFrom1[T0, T1, R any](
	ob1 Observable[T1],
	mapping func(v0 T0, v1 T1) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom2(source, ob1, mapping)
		},
	)
}

func withLatestFrom2[T1, T2, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s withLatestFromState2[T1, T2]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { withLatestFromEmit2(o, n, mapping, &s, &s.v1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { withLatestFromEmit2(o, n, mapping, &s, &s.v2, 2) })
	}
}

type withLatestFromState2[T1, T2 any] struct {
	mu sync.Mutex

	nbits uint8

	v1 T1
	v2 T2
}

func withLatestFromEmit2[T1, T2, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2) R,
	s *withLatestFromState2[T1, T2],
	v *X,
	bit uint8,
) {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		*v = n.Value
		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits && bit == 1 {
			v := Try21(mapping, s.v1, s.v2, s.mu.Unlock)
			s.mu.Unlock()
			o.Next(v)
			return
		}

		s.mu.Unlock()

	case KindError:
		o.Error(n.Error)

	case KindComplete:
		if bit == 1 {
			o.Complete()
		}
	}
}
