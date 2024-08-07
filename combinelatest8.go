package rx

import "sync"

// CombineLatest8 combines multiple Observables to create an Observable
// that emits mappings of the latest values emitted by each of its input
// Observables.
func CombineLatest8[T1, T2, T3, T4, T5, T6, T7, T8, R any](
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
		c, o = Serialize(c, o)

		var s combineLatestState8[T1, T2, T3, T4, T5, T6, T7, T8]

		_ = true &&
			ob1.satcc(c, func(n Notification[T1]) { combineLatestEmit8(o, n, mapping, &s, &s.v1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { combineLatestEmit8(o, n, mapping, &s, &s.v2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { combineLatestEmit8(o, n, mapping, &s, &s.v3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { combineLatestEmit8(o, n, mapping, &s, &s.v4, 8) }) &&
			ob5.satcc(c, func(n Notification[T5]) { combineLatestEmit8(o, n, mapping, &s, &s.v5, 16) }) &&
			ob6.satcc(c, func(n Notification[T6]) { combineLatestEmit8(o, n, mapping, &s, &s.v6, 32) }) &&
			ob7.satcc(c, func(n Notification[T7]) { combineLatestEmit8(o, n, mapping, &s, &s.v7, 64) }) &&
			ob8.satcc(c, func(n Notification[T8]) { combineLatestEmit8(o, n, mapping, &s, &s.v8, 128) })
	}
}

type combineLatestState8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	mu sync.Mutex

	nbits, cbits uint8

	v1 T1
	v2 T2
	v3 T3
	v4 T4
	v5 T5
	v6 T6
	v7 T7
	v8 T8
}

func combineLatestEmit8[T1, T2, T3, T4, T5, T6, T7, T8, R, X any](
	o Observer[R],
	n Notification[X],
	mapping func(T1, T2, T3, T4, T5, T6, T7, T8) R,
	s *combineLatestState8[T1, T2, T3, T4, T5, T6, T7, T8],
	v *X,
	bit uint8,
) {
	const FullBits = 255

	switch n.Kind {
	case KindNext:
		s.mu.Lock()
		*v = n.Value
		nbits := s.nbits
		nbits |= bit
		s.nbits = nbits

		if nbits == FullBits {
			v := Try81(mapping, s.v1, s.v2, s.v3, s.v4, s.v5, s.v6, s.v7, s.v8, s.mu.Unlock)
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
