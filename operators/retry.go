package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/misc"
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
		if t.HasValue || !t.HasError || count == 0 {
			sink(t)
			return
		}
		if count > 0 {
			count--
		}
		avoidRecursion.Do(subscribe)
	}

	avoidRecursion.Do(subscribe)
}

// Retry creates an Observable that mirrors the source Observable with one
// exception: when the source emits an error, this operator will resubscribe
// to the source for a maximum of count resubscriptions rather than propagating
// this error.
func Retry(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count == 0 {
			return source
		}
		return retryObservable{source, count}.Subscribe
	}
}

// RetryForever creates an Observable that mirrors the source Observable with
// one exception: when the source emits an error, this operator will always
// resubscribe to the source rather than propagating this error.
func RetryForever() rx.Operator {
	return Retry(-1)
}
