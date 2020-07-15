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
		subject        *rx.Subject
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
			if subject == nil {
				subject = rx.NewSubject()
				obs := obs.Notifier(subject.Observable)
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
			subject.Next(nil)
		})
	}

	avoidRecursion.Do(subscribe)
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// the exception of COMPLETE emission. If the source Observable completes,
// this operator will emit nil to the Observable returned from notifier. If
// that Observable emits a value, this operator will resubscribe to the source
// Observable. Otherwise, this operator will emit a COMPLETE on the child
// subscription.
func RepeatWhen(notifier rx.Operator) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return repeatWhenObservable{source, notifier}.Subscribe
	}
}
