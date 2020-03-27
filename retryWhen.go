package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type retryWhenObservable struct {
	Source   Observable
	Notifier func(Observable) Observable
}

func (obs retryWhenObservable) Subscribe(ctx context.Context, sink Observer) {
	sink = Mutex(sink)

	var (
		activeCount    = atomic.Uint32(2)
		lastError      error
		subject        Subject
		createSubject  func() Subject
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
		obs.Source.Subscribe(sourceCtx, func(t Notification) {
			if x, ok := <-cx; ok {
				switch {
				case t.HasValue:
					sink(t)
					cx <- x
				case t.HasError:
					lastError = t.Error
					activeCount := activeCount.Sub(1)
					close(cx)
					if activeCount == 0 {
						sink(t)
						break
					}
					if subject.Observer == nil {
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

	createSubject = func() Subject {
		subject := NewSubject()
		obs := obs.Notifier(subject.Observable)
		obs.Subscribe(ctx, func(t Notification) {
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
					sink.Error(lastError)
				}
			}
		})
		return subject
	}

	avoidRecursive.Do(subscribe)
}

// RetryWhen creates an Observable that mirrors the source Observable with
// the exception of ERROR emission. If the source Observable errors, this
// operator will emit the error to the Observable returned from notifier.
// If that Observable emits a value, this operator will resubscribe to the
// source Observable. Otherwise, this operator will emit the last error on
// the child subscription.
func (Operators) RetryWhen(notifier func(Observable) Observable) Operator {
	return func(source Observable) Observable {
		obs := retryWhenObservable{source, notifier}
		return Create(obs.Subscribe)
	}
}
