package rx

import "sync"

// CombineLatest2 combines multiple Observables to create an [Observable]
// that emits mappings of the latest values emitted by each of the input
// Observables.
func CombineLatest2[T1, T2, R any](
	ob1 Observable[T1],
	ob2 Observable[T2],
	mapping func(v1 T1, v2 T2) R,
) Observable[R] {
	return func(c Context, o Observer[R]) {
		c, o = Synchronize(c, o)

		var s combineLatestState2[T1, T2]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { combineLatestEmit2(o, n, mapping, &s, &s.v1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { combineLatestEmit2(o, n, mapping, &s, &s.v2, 2) })
	}
}

type combineLatestState2[T1, T2 any] struct {
	mu sync.Mutex

	nbits, cbits uint8

	v1 T1
	v2 T2
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
		s.mu.Lock()
		*v = n.Value
		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits {
			v := Try21(mapping, s.v1, s.v2, s.mu.Unlock)
			s.mu.Unlock()
			o.Next(v)
			return
		}

		s.mu.Unlock()

	case KindComplete:
		s.mu.Lock()
		cbits := s.cbits
		cbits |= bit
		s.cbits = cbits
		s.mu.Unlock()

		if cbits == FullBits {
			o.Complete()
		}

	case KindError:
		o.Error(n.Error)

	case KindStop:
		o.Stop(n.Error)
	}
}
