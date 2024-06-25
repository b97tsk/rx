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
			ob1.satcc(c, func(n Notification[T1]) { combineLatestEmit8(o, n, mapping, &s, &s.V1, 1) }) &&
			ob2.satcc(c, func(n Notification[T2]) { combineLatestEmit8(o, n, mapping, &s, &s.V2, 2) }) &&
			ob3.satcc(c, func(n Notification[T3]) { combineLatestEmit8(o, n, mapping, &s, &s.V3, 4) }) &&
			ob4.satcc(c, func(n Notification[T4]) { combineLatestEmit8(o, n, mapping, &s, &s.V4, 8) }) &&
			ob5.satcc(c, func(n Notification[T5]) { combineLatestEmit8(o, n, mapping, &s, &s.V5, 16) }) &&
			ob6.satcc(c, func(n Notification[T6]) { combineLatestEmit8(o, n, mapping, &s, &s.V6, 32) }) &&
			ob7.satcc(c, func(n Notification[T7]) { combineLatestEmit8(o, n, mapping, &s, &s.V7, 64) }) &&
			ob8.satcc(c, func(n Notification[T8]) { combineLatestEmit8(o, n, mapping, &s, &s.V8, 128) })
	}
}

type combineLatestState8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	sync.Mutex

	NBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
	V8 T8
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
		s.Lock()
		*v = n.Value
		nbits := s.NBits
		nbits |= bit
		s.NBits = nbits

		if nbits == FullBits {
			v := Try81(mapping, s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7, s.V8, s.Unlock)
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
