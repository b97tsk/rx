package rx

import "sync"

// WithLatestFrom2 combines the source with 2 other Observables to create
// an Observable that emits mappings of the latest values emitted by each
// Observable, only when the source emits.
func WithLatestFrom2[T0, T1, T2, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	mapping func(v0 T0, v1 T1, v2 T2) R,
) Operator[T0, R] {
	return NewOperator(
		func(source Observable[T0]) Observable[R] {
			return withLatestFrom3(source, ob1, ob2, mapping)
		},
	)
}

func withLatestFrom3[T1, T2, T3, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	mapping func(v1 T1, v2 T2, v3 T3) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s withLatestFromState3[T1, T2, T3]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { withLatestFromEmit3(o, n, mapping, &s, &s.V1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { withLatestFromEmit3(o, n, mapping, &s, &s.V2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { withLatestFromEmit3(o, n, mapping, &s, &s.V3, 4) })
	}
}

type withLatestFromState3[T1, T2, T3 any] struct {
	sync.Mutex

	NBits uint8

	V1 T1
	V2 T2
	V3 T3
}

func withLatestFromEmit3[T1, T2, T3, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3) R,
	s *withLatestFromState3[T1, T2, T3],
	v *X,
	bit uint8,
) {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		s.Lock()
		*v = n.Value
		nbits := s.NBits
		nbits |= bit
		s.NBits = nbits

		if nbits == FullBits && bit == 1 {
			v := Try31(mapping, s.V1, s.V2, s.V3, s.Unlock)
			s.Unlock()
			o.Next(v)
			return
		}

		s.Unlock()

	case KindError:
		o.Error(n.Error)

	case KindComplete:
		if bit == 1 {
			o.Complete()
		}
	}
}
