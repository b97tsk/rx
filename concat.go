package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/norec"
)

// Concat creates an output Observable which sequentially emits all values
// from given Observable and then moves on to the next.
//
// Concat concatenates multiple Observables together by sequentially emitting
// their values, one Observable after the other.
func Concat(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}

	return concatObservable(observables).Subscribe
}

type concatObservable []Observable

func (observables concatObservable) Subscribe(ctx context.Context, sink Observer) {
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
		if t.HasValue || t.HasError {
			sink(t)
			return
		}
		subscribeToNext()
	}

	subscribeToNext()
}
