package rx

import (
	"context"

	"github.com/b97tsk/rx/x/misc"
)

// OnErrorResumeNext creates an Observable that concatenates the provided
// Observables. It's like Concat, but errors are ignored.
func OnErrorResumeNext(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return Create(
		func(ctx context.Context, sink Observer) {
			var (
				observer       Observer
				avoidRecursion misc.AvoidRecursion
			)

			remainder := observables
			subscribe := func() {
				if ctx.Err() != nil {
					return
				}
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
				avoidRecursion.Do(subscribe)
			}

			avoidRecursion.Do(subscribe)
		},
	)
}
