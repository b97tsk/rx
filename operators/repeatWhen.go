package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/norec"
)

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

type repeatWhenObservable struct {
	Source   rx.Observable
	Notifier rx.Operator
}

func (obs repeatWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	repeatWorking := atomic.Uint32(1)
	sourceWorking := atomic.Uint32(1)

	var repeatSignal rx.Observer

	var subscribeToSource func()

	subscribeToSource = norec.Wrap(func() {
		obs.Source.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue || t.HasError {
				sink(t)

				return
			}

			sourceWorking.Store(0)

			if repeatWorking.Equals(0) {
				sink(t)

				return
			}

			if repeatSignal == nil {
				d := rx.Unicast()

				repeatSignal = d.Observer

				obs1 := obs.Notifier(d.Observable)

				obs1.Subscribe(ctx, func(t rx.Notification) {
					switch {
					case t.HasValue:
						if sourceWorking.Cas(0, 1) {
							subscribeToSource()
						}

					case t.HasError:
						sink(t)

					default:
						repeatWorking.Store(0)

						if sourceWorking.Equals(0) {
							sink(t)
						}
					}
				})
			}

			repeatSignal.Next(nil)
		})
	})

	subscribeToSource()
}
