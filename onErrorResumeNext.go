package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/norec"
)

// OnErrorResumeNext creates an Observable that concatenates multiple
// Observables together by sequentially emitting their values, one Observable
// after the other. It's like Concat, but errors are ignored.
func OnErrorResumeNext(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}

	return onErrorResumeNextObservable(observables).Subscribe
}

type onErrorResumeNextObservable []Observable

func (observables onErrorResumeNextObservable) Subscribe(ctx context.Context, sink Observer) {
	var observer Observer

	remainder := observables

	subscribeToNext := norec.Wrap(func() {
		if len(remainder) == 0 {
			sink.Complete()

			return
		}

		if ctx.Err() != nil {
			return
		}

		source := remainder[0]
		remainder = remainder[1:]

		source.Subscribe(ctx, observer)
	})

	observer = func(t Notification) {
		if t.HasValue {
			sink(t)

			return
		}

		subscribeToNext()
	}

	subscribeToNext()
}
