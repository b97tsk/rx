package rx

import "sync"

// CombineLatest3 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest3[T1, T2, T3, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	ob3 Observable[T3],
	mapping func(v1 T1, v2 T2, v3 T3) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s combineLatestState3[T1, T2, T3]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { combineLatestEmit3(o, n, mapping, &s, &s.v1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { combineLatestEmit3(o, n, mapping, &s, &s.v2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { combineLatestEmit3(o, n, mapping, &s, &s.v3, 4) })
	}
}

type combineLatestState3[T1, T2, T3 any] struct {
	mu sync.Mutex

	nbits, cbits uint8

	v1 T1
	v2 T2
	v3 T3
}

func combineLatestEmit3[T1, T2, T3, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3) R,
	s *combineLatestState3[T1, T2, T3],
	v *X,
	bit uint8,
) {
	const FullBits = 7

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		*v = n.Value
		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits {
			v := Try31(mapping, s.v1, s.v2, s.v3, s.mu.Unlock)
			s.mu.Unlock()
			o.Next(v)
			return
		}

		s.mu.Unlock()

	case KindError:
		o.Error(n.Error)

	case KindComplete:
		s.mu.Lock()
		cbits := s.cbits
		cbits |= bit
		s.cbits = cbits
		s.mu.Unlock()

		if cbits == FullBits {
			o.Complete()
		}
	}
}
