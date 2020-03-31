package rx

import (
	"context"

	"github.com/b97tsk/rx/x/misc"
)

type onErrorResumeNextObservable struct {
	Observables []Observable
}

func (obs onErrorResumeNextObservable) Subscribe(ctx context.Context, sink Observer) {
	var (
		observer       Observer
		avoidRecursive misc.AvoidRecursive
	)

	remainder := obs.Observables
	subscribe := func() {
		source := remainder[0]
		remainder = remainder[1:]
		source.Subscribe(ctx, observer)
	}

	observer = func(t Notification) {
		if t.HasValue {
			sink(t)
			return
		}
		if len(remainder) == 0 {
			sink.Complete()
			return
		}
		avoidRecursive.Do(subscribe)
	}

	avoidRecursive.Do(subscribe)
}

// OnErrorResumeNext creates an Observable that concatenates the provided
// Observables. It's like Concat, but errors are ignored.
func OnErrorResumeNext(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := onErrorResumeNextObservable{observables}
	return Create(obs.Subscribe)
}
