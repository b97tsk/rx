package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/atomic"
)

type repeatWhenObservable struct {
	Source   rx.Observable
	Notifier func(rx.Observable) rx.Observable
}

func (obs repeatWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	var (
		activeCount    = atomic.Uint32(2)
		subject        rx.Subject
		createSubject  func() rx.Subject
		avoidRecursive avoidRecursiveCalls
	)

	type X struct{}
	var cxCurrent chan X

	var (
		sourceCtx    context.Context
		sourceCancel context.CancelFunc
	)

	subscribe := func() {
		cx := make(chan X, 1)
		cx <- X{}
		cxCurrent = cx
		sourceCtx, sourceCancel = context.WithCancel(ctx)
		obs.Source.Subscribe(sourceCtx, func(t rx.Notification) {
			if x, ok := <-cx; ok {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
					cx <- x
				default:
					activeCount := activeCount.Sub(1)
					close(cx)
					if activeCount == 0 {
						sink(t)
						break
					}
					if subject.Observer == nil {
						subject = createSubject()
					}
					subject.Next(nil)
				}
			}
		})
	}

	createSubject = func() rx.Subject {
		subject := rx.NewSubject()
		obs := obs.Notifier(subject.Observable)
		obs.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				sourceCancel()
				if _, ok := <-cxCurrent; ok {
					close(cxCurrent)
				} else {
					activeCount.Add(1)
				}
				avoidRecursive.Do(subscribe)

			case t.HasError:
				sink(t)

			default:
				if activeCount.Sub(1) == 0 {
					sink(t)
				}
			}
		})
		return subject
	}

	avoidRecursive.Do(subscribe)
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// the exception of COMPLETE emission. If the source Observable completes,
// this operator will emit nil to the Observable returned from notifier. If
// that Observable emits a value, this operator will resubscribe to the source
// Observable. Otherwise, this operator will emit a COMPLETE on the child
// subscription.
func RepeatWhen(notifier func(rx.Observable) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := repeatWhenObservable{source, notifier}
		return rx.Create(obs.Subscribe)
	}
}
