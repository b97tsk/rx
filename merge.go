package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type mergeObservable []Observable

func (observables mergeObservable) Subscribe(ctx context.Context, sink Observer) {
	sink = Mutex(sink)
	active := atomic.Uint32(len(observables))
	observer := func(t Notification) {
		if t.HasValue || t.HasError || active.Sub(1) == 0 {
			sink(t)
		}
	}
	for _, obs := range observables {
		go obs.Subscribe(ctx, observer)
	}
}

// Merge creates an output Observable which concurrently emits all values from
// every given input Observable.
//
// Merge flattens multiple Observables together by blending their values into
// one Observable.
func Merge(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return Create(mergeObservable(observables).Subscribe)
}
