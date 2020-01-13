package rx

import (
	"context"
)

type onErrorResumeNextOperator struct {
	Observables []Observable
}

func (op onErrorResumeNextOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		observer       Observer
		avoidRecursive avoidRecursiveCalls
	)

	remainder := op.Observables
	subscribe := func() {
		source, remainder = remainder[0], remainder[1:]
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

	return ctx, cancel
}

// OnErrorResumeNext creates an Observable that concatenates the provided
// Observables. It's like Concat, but errors are ignored.
func OnErrorResumeNext(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	op := onErrorResumeNextOperator{observables}
	return Empty().Lift(op.Call)
}
