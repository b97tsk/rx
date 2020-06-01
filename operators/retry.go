package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/misc"
)

type retryObservable struct {
	Source rx.Observable
	Count  int
}

func (obs retryObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		count          = obs.Count
		observer       rx.Observer
		avoidRecursion misc.AvoidRecursion
	)

	subscribe := func() {
		obs.Source.Subscribe(ctx, observer)
	}

	observer = func(t rx.Notification) {
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
				avoidRecursion.Do(subscribe)
			}
		default:
			sink(t)
		}
	}

	avoidRecursion.Do(subscribe)
}

// Retry creates an Observable that mirrors the source Observable with the
// exception of ERROR emission. If the source Observable errors, this
// operator will resubscribe to the source Observable for a maximum of count
// resubscriptions rather than propagating the ERROR emission.
func Retry(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count == 0 {
			return source
		}
		obs := retryObservable{source, count}
		return rx.Create(obs.Subscribe)
	}
}
