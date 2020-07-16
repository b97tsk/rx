package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
	"github.com/b97tsk/rx/subject"
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
		subject1       *subject.Subject
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
			if subject1 == nil {
				subject1 = subject.NewSubject()
				obs := obs.Notifier(subject1.Observable)
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
			subject1.Next(t.Error)
		})
	}

	avoidRecursion.Do(subscribe)
}

// RetryWhen creates an Observable that mirrors the source Observable with
// the exception of ERROR emission. If the source Observable errors, this
// operator will emit the error to the Observable returned from notifier.
// If that Observable emits a value, this operator will resubscribe to the
// source Observable. Otherwise, this operator will emit the last error on
// the child subscription.
func RetryWhen(notifier rx.Operator) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return retryWhenObservable{source, notifier}.Subscribe
	}
}
