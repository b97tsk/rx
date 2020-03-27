package rx

import (
	"context"
	"time"
)

// A TimeoutConfigure is a configure for Timeout.
type TimeoutConfigure struct {
	Duration   time.Duration
	Observable Observable
}

// Use creates an Operator from this configure.
func (configure TimeoutConfigure) Use() Operator {
	return func(source Observable) Observable {
		obs := timeoutObservable{source, configure}
		return Create(obs.Subscribe)
	}
}

type timeoutObservable struct {
	Source Observable
	TimeoutConfigure
}

func (obs timeoutObservable) Subscribe(ctx context.Context, sink Observer) {
	childCtx, childCancel := context.WithCancel(ctx)

	type X struct{}
	cx := make(chan X, 1)
	cx <- X{}

	var scheduleCancel context.CancelFunc

	doSchedule := func() {
		if scheduleCancel != nil {
			scheduleCancel()
		}
		_, scheduleCancel = scheduleOnce(childCtx, obs.Duration, func() {
			if _, ok := <-cx; ok {
				close(cx)
				childCancel()
				obs.Observable.Subscribe(ctx, sink)
			}
		})
	}

	doSchedule()

	obs.Source.Subscribe(childCtx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sink(t)
				cx <- x
				doSchedule()
			default:
				close(cx)
				sink(t)
			}
		}
	})
}

// Timeout creates an Observable that mirrors the source Observable or notify
// of an ErrTimeout if the source does not emit a value in given time span.
func (Operators) Timeout(timeout time.Duration) Operator {
	return TimeoutConfigure{timeout, Throw(ErrTimeout)}.Use()
}
