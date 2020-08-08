package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
)

type repeatWhenObservable struct {
	Source   rx.Observable
	Notifier rx.Operator
}

func (obs repeatWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	var (
		repeatActive   = atomic.Uint32(1)
		sourceActive   = atomic.Uint32(1)
		repeatSignal   rx.Observer
		subscribe      func()
		avoidRecursion misc.AvoidRecursion
	)

	subscribe = func() {
		obs.Source.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue || t.HasError {
				sink(t)
				return
			}
			sourceActive.Store(0)
			if repeatActive.Equals(0) {
				sink(t)
				return
			}
			if repeatSignal == nil {
				d := rx.Unicast()
				repeatSignal = d.Observer
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
						repeatActive.Store(0)
						if sourceActive.Equals(0) {
							sink(t)
						}
					}
				})
			}
			repeatSignal.Next(nil)
		})
	}

	avoidRecursion.Do(subscribe)
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// one exception: when the source completes, this operator will emit nil to
// the Observable returned by the notifier. If that Observable emits a value,
// this operator will resubscribe to the source; otherwise, this operator
// will cause the child subscription to complete.
func RepeatWhen(notifier rx.Operator) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return repeatWhenObservable{source, notifier}.Subscribe
	}
}
