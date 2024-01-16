package rx

import "context"

// Zip3 combines multiple Observables to create an Observable that emits
// projection of values of each of its input Observables.
//
// Zip3 pulls values from each input Observable one by one, it only buffers
// one value for each input Observable.
func Zip3[T1, T2, T3, R any](
	obs1 Observable[T1],
	obs2 Observable[T2],
	obs3 Observable[T3],
	proj func(v1 T1, v2 T2, v3 T3) R,
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

		chan1 := make(chan Notification[T1], 1)
		chan2 := make(chan Notification[T2], 1)
		chan3 := make(chan Notification[T3], 1)

		wg.Go(func() { obs1.Subscribe(ctx, channelObserver(chan1, noop)) })
		wg.Go(func() { obs2.Subscribe(ctx, channelObserver(chan2, noop)) })
		wg.Go(func() { obs3.Subscribe(ctx, channelObserver(chan3, noop)) })

		wg.Go(func() {
			for {
			Again1:
				n1 := <-chan1
				switch n1.Kind {
				case KindNext:
				case KindError:
					sink.Error(n1.Error)
					return
				case KindComplete:
					sink.Complete()
					return
				default:
					goto Again1
				}
			Again2:
				n2 := <-chan2
				switch n2.Kind {
				case KindNext:
				case KindError:
					sink.Error(n2.Error)
					return
				case KindComplete:
					sink.Complete()
					return
				default:
					goto Again2
				}
			Again3:
				n3 := <-chan3
				switch n3.Kind {
				case KindNext:
				case KindError:
					sink.Error(n3.Error)
					return
				case KindComplete:
					sink.Complete()
					return
				default:
					goto Again3
				}
				sink.Next(proj(n1.Value, n2.Value, n3.Value))
			}
		})
	}
}
