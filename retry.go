package rx

import (
	"context"
)

type retryObservable struct {
	Source Observable
	Count  int
}

func (obs retryObservable) Subscribe(ctx context.Context, sink Observer) {
	var (
		count          = obs.Count
		observer       Observer
		avoidRecursive avoidRecursiveCalls
	)

	subscribe := func() {
		obs.Source.Subscribe(ctx, observer)
	}

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			if count == 0 {
				sink(t)
			} else {
				if count > 0 {
					count--
				}
				avoidRecursive.Do(subscribe)
			}
		default:
			sink(t)
		}
	}

	avoidRecursive.Do(subscribe)
}

// Retry creates an Observable that mirrors the source Observable with the
// exception of ERROR emission. If the source Observable errors, this
// operator will resubscribe to the source Observable for a maximum of count
// resubscriptions rather than propagating the ERROR emission.
func (Operators) Retry(count int) Operator {
	return func(source Observable) Observable {
		if count == 0 {
			return source
		}
		obs := retryObservable{source, count}
		return Create(obs.Subscribe)
	}
}
