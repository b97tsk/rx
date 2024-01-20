package rx

import "context"

// CombineLatest7 combines multiple Observables to create an Observable that
// emits projection of latest values of each of its input Observables.
func CombineLatest7[T1, T2, T3, T4, T5, T6, T7, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	obs4 Observable[T4],
	obs5 Observable[T5],
	obs6 Observable[T6],
	obs7 Observable[T7],
	proj func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
) Observable[R] {
	if proj == nil {
		panic("proj == nil")
	}

	return func(ctx context.Context, sink Observer[R]) {
		wg := WaitGroupFromContext(ctx)
		ctx, cancel := context.WithCancel(ctx)

		noop := make(chan struct{})

		sink = sink.OnLastNotification(func() {
			cancel()
			close(noop)
		})

		chan1 := make(chan Notification[T1])
		chan2 := make(chan Notification[T2])
		chan3 := make(chan Notification[T3])
		chan4 := make(chan Notification[T4])
		chan5 := make(chan Notification[T5])
		chan6 := make(chan Notification[T6])
		chan7 := make(chan Notification[T7])

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, channelObserver(chan3, noop)) })
		wg.Go(func() { obs4.Subscribe(ctx, channelObserver(chan4, noop)) })
		wg.Go(func() { obs5.Subscribe(ctx, channelObserver(chan5, noop)) })
		wg.Go(func() { obs6.Subscribe(ctx, channelObserver(chan6, noop)) })
		wg.Go(func() { obs7.Subscribe(ctx, channelObserver(chan7, noop)) })

		wg.Go(func() {
			var s combineLatestState7[T1, T2, T3, T4, T5, T6, T7]

			cont := true

			for cont {
				select {
				case n := <-chan1:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V1, 1)
				case n := <-chan2:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V2, 2)
				case n := <-chan3:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V3, 4)
				case n := <-chan4:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V4, 8)
				case n := <-chan5:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V5, 16)
				case n := <-chan6:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V6, 32)
				case n := <-chan7:
					cont = combineLatestSink7(n, sink, proj, &s, &s.V7, 64)
				}
			}
		})
	}
}

type combineLatestState7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	NBits, CBits uint8

	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
}

func combineLatestSink7[T1, T2, T3, T4, T5, T6, T7, R, X any](
	n Notification[X],
	sink Observer[R],
	proj func(T1, T2, T3, T4, T5, T6, T7) R,
	s *combineLatestState7[T1, T2, T3, T4, T5, T6, T7],
	v *X,
	bit uint8,
) bool {
	const FullBits = 127

	switch n.Kind {
	case KindNext:
		*v = n.Value

		if s.NBits |= bit; s.NBits == FullBits {
			sink.Next(proj(s.V1, s.V2, s.V3, s.V4, s.V5, s.V6, s.V7))
		}

	case KindError:
		sink.Error(n.Error)
		return false

	case KindComplete:
		if s.CBits |= bit; s.CBits == FullBits {
			sink.Complete()
			return false
		}
	}

	return true
}
