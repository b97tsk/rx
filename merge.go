package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/atomic"
)

// Merge creates an output Observable which concurrently emits all values from
// every given input Observable.
//
// Merge flattens multiple Observables together by blending their values into
// one Observable.
func Merge(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}

	return mergeObservable(observables).Subscribe
}

type mergeObservable []Observable

func (observables mergeObservable) Subscribe(ctx context.Context, sink Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	active := atomic.Uint32(uint32(len(observables)))

	observer := func(t Notification) {
		if t.HasValue || t.HasError || active.Sub(1) == 0 {
			sink(t)
		}
	}

	for _, obs := range observables {
		go obs.Subscribe(ctx, observer)
	}
}
