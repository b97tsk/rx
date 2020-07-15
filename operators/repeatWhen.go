package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
)

type repeatWhenObservable struct {
	Source   rx.Observable
	Notifier func(rx.Observable) rx.Observable
}

func (obs repeatWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = sink.Mutex()

	var (
		active         = atomic.Uint32(2)
		subject        *rx.Subject
		createSubject  func() *rx.Subject
		avoidRecursion misc.AvoidRecursion
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
				if t.HasValue || t.HasError {
					sink(t)
					cx <- x
					return
				}
				active := active.Sub(1)
				close(cx)
				if active == 0 {
					sink(t)
					return
				}
				if subject == nil {
					subject = createSubject()
				}
				subject.Next(nil)
			}
		})
	}

	createSubject = func() *rx.Subject {
		subject := rx.NewSubject()
		obs := obs.Notifier(subject.Observable)
		obs.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				sourceCancel()
				if _, ok := <-cxCurrent; ok {
					close(cxCurrent)
				} else {
					active.Add(1)
				}
				avoidRecursion.Do(subscribe)

			case t.HasError:
				sink(t)

			default:
				if active.Sub(1) == 0 {
					sink(t)
				}
			}
		})
		return subject
	}

	avoidRecursion.Do(subscribe)
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
