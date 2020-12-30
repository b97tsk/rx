package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/norec"
)

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

type retryWhenObservable struct {
	Source   rx.Observable
	Notifier rx.Operator
}

func (obs retryWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	retryWorking := atomic.FromBool(true)
	sourceWorking := atomic.FromBool(true)

	var lastError error

	var retrySignal rx.Observer

	var subscribeToSource func()

	subscribeToSource = norec.Wrap(func() {
		obs.Source.Subscribe(ctx, func(t rx.Notification) {
			if !t.HasError {
				sink(t)

				return
			}

			lastError = t.Error

			sourceWorking.Store(false)

			if retryWorking.False() {
				sink(t)

				return
			}

			if retrySignal == nil {
				d := rx.Unicast()

				retrySignal = d.Observer

				obs1 := obs.Notifier(d.Observable)

				obs1.Subscribe(ctx, func(t rx.Notification) {
					switch {
					case t.HasValue:
						if sourceWorking.Cas(false, true) {
							subscribeToSource()
						}

					case t.HasError:
						sink(t)

					default:
						retryWorking.Store(false)

						if sourceWorking.False() {
							sink.Error(lastError)
						}
					}
				})
			}

			retrySignal.Next(t.Error)
		})
	})

	subscribeToSource()
}
