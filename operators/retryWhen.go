package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/misc"
)

type retryWhenObservable struct {
	Source   rx.Observable
	Notifier func(rx.Observable) rx.Observable
}

func (obs retryWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	var (
		err            error
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
				switch {
				case t.HasValue:
					sink(t)
					cx <- x
				case t.HasError:
					err = t.Error
					active := active.Sub(1)
					close(cx)
					if active == 0 {
						sink(t)
						break
					}
					if subject == nil {
						subject = createSubject()
					}
					subject.Next(t.Error)
				default:
					sink(t)
					cx <- x
				}
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
					sink.Error(err)
				}
			}
		})
		return subject
	}

	avoidRecursion.Do(subscribe)
}

// RetryWhen creates an Observable that mirrors the source Observable with
// the exception of ERROR emission. If the source Observable errors, this
// operator will emit the error to the Observable returned from notifier.
// If that Observable emits a value, this operator will resubscribe to the
// source Observable. Otherwise, this operator will emit the last error on
// the child subscription.
func RetryWhen(notifier func(rx.Observable) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := retryWhenObservable{source, notifier}
		return rx.Create(obs.Subscribe)
	}
}
