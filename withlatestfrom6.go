package rx

import "sync"

// WithLatestFrom6 combines the source with 6 other Observables to create
// an [Observable] that emits mappings of the latest values emitted by each
// [Observable], only when the source emits.
func WithLatestFrom6[T0, T1, T2, T3, T4, T5, T6, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	ob6 Observable[T6],
	mapping func(v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom7(source, ob1, ob2, ob3, ob4, ob5, ob6, mapping)
		},
	)
}

func withLatestFrom7[T1, T2, T3, T4, T5, T6, T7, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	ob4 Observable[T4],
	ob5 Observable[T5],
	ob6 Observable[T6],
	ob7 Observable[T7],
	mapping func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s withLatestFromState7[T1, T2, T3, T4, T5, T6, T7]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { withLatestFromEmit7(o, n, mapping, &s, &s.v1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { withLatestFromEmit7(o, n, mapping, &s, &s.v2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { withLatestFromEmit7(o, n, mapping, &s, &s.v3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { withLatestFromEmit7(o, n, mapping, &s, &s.v4, 8) }) &&
			ob5.satcc(c, func(n Notification[T5]) { withLatestFromEmit7(o, n, mapping, &s, &s.v5, 16) }) &&
			ob6.satcc(c, func(n Notification[T6]) { withLatestFromEmit7(o, n, mapping, &s, &s.v6, 32) }) &&
			ob7.satcc(c, func(n Notification[T7]) { withLatestFromEmit7(o, n, mapping, &s, &s.v7, 64) })
	}
}

type withLatestFromState7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	mu sync.Mutex

	nbits uint8

	v1 T1
	v2 T2
	v3 T3
	v4 T4
	v5 T5
	v6 T6
	v7 T7
}

func withLatestFromEmit7[T1, T2, T3, T4, T5, T6, T7, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5, T6, T7) R,
	s *withLatestFromState7[T1, T2, T3, T4, T5, T6, T7],
	v *X,
	bit uint8,
) {
	const FullBits = 127

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		*v = n.Value
		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits && bit == 1 {
			v := Try71(mapping, s.v1, s.v2, s.v3, s.v4, s.v5, s.v6, s.v7, s.mu.Unlock)
			s.mu.Unlock()
			o.Next(v)
			return
		}

		s.mu.Unlock()

	case KindComplete:
		if bit == 1 {
			o.Complete()
		}

	case KindError:
		o.Error(n.Error)

	case KindStop:
		o.Stop(n.Error)
	}
}
