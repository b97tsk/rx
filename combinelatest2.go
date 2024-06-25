package rx

import "sync"

// CombineLatest2 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest2[T1, T2, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Serialize(c, o)

		var s combineLatestState2[T1, T2]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { combineLatestEmit2(o, n, mapping, &s, &s.V1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { combineLatestEmit2(o, n, mapping, &s, &s.V2, 2) })
	}
}

type combineLatestState2[T1, T2 any] struct {
	sync.Mutex

	NBits, CBits uint8

	V1 T1
	V2 T2
}

func combineLatestEmit2[T1, T2, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2) R,
	s *combineLatestState2[T1, T2],
	v *X,
	bit uint8,
) {
	const FullBits = 3

	switch n.Kind {
	case KindNext:
		s.Lock()
		*v = n.Value
		nbits := s.NBits
		nbits |= bit
		s.NBits = nbits

		if nbits == FullBits {
			v := Try21(mapping, s.V1, s.V2, s.Unlock)
			s.Unlock()
			o.Next(v)
			return
		}

		s.Unlock()

	case KindError:
		o.Error(n.Error)

	case KindComplete:
		s.Lock()
		cbits := s.CBits
		cbits |= bit
		s.CBits = cbits
		s.Unlock()

		if cbits == FullBits {
			o.Complete()
		}
	}
}
