package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
)

type retryWhenObservable struct {
	Source   rx.Observable
	Notifier rx.Operator
}

func (obs retryWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	var (
		err            error
		retryActive    = atomic.Uint32(1)
		sourceActive   = atomic.Uint32(1)
		retrySignal    rx.Observer
		subscribe      func()
		avoidRecursion misc.AvoidRecursion
	)

	subscribe = func() {
		obs.Source.Subscribe(ctx, func(t rx.Notification) {
			if !t.HasError {
				sink(t)
				return
			}
			err = t.Error
			sourceActive.Store(0)
			if retryActive.Equals(0) {
				sink(t)
				return
			}
			if retrySignal == nil {
				d := rx.Unicast()
				retrySignal = d.Observer
				obs := obs.Notifier(d.Observable)
				obs.Subscribe(ctx, func(t rx.Notification) {
					switch {
					case t.HasValue:
						if sourceActive.Cas(0, 1) {
							avoidRecursion.Do(subscribe)
						}
					case t.HasError:
						sink(t)
					default:
						retryActive.Store(0)
						if sourceActive.Equals(0) {
							sink.Error(err)
						}
					}
				})
			}
			retrySignal.Next(t.Error)
		})
	}

	avoidRecursion.Do(subscribe)
}

// RetryWhen creates an Observable that mirrors the source Observable with
// one exception: when the source emits an error, this operator will emit
// this error to the Observable returned by the notifier. If that Observable
// emits a value, this operator will resubscribe to the source; otherwise,
// this operator will emit the last error on the child subscription.
func RetryWhen(notifier rx.Operator) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return retryWhenObservable{source, notifier}.Subscribe
	}
}
