package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

// A TimeoutConfigure is a configure for Timeout.
type TimeoutConfigure struct {
	Duration   time.Duration
	Observable rx.Observable
}

// Use creates an Operator from this configure.
func (configure TimeoutConfigure) Use() rx.Operator {
	if configure.Observable == nil {
		configure.Observable = rx.Throw(rx.ErrTimeout)
	}
	return func(source rx.Observable) rx.Observable {
		obs := timeoutObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type timeoutObservable struct {
	Source rx.Observable
	TimeoutConfigure
}

func (obs timeoutObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	childCtx, childCancel := context.WithCancel(ctx)

	type X struct{}
	cx := make(chan X, 1)
	cx <- X{}

	var scheduleCancel context.CancelFunc

	obsTimer := rx.Timer(obs.Duration)
	doSchedule := func() {
		if scheduleCancel != nil {
			scheduleCancel()
		}
		_, scheduleCancel = obsTimer.Subscribe(childCtx, func(t rx.Notification) {
			if t.HasValue {
				return
			}
			if _, ok := <-cx; ok {
				close(cx)
				childCancel()
				obs.Observable.Subscribe(ctx, sink)
			}
		})
	}

	doSchedule()

	obs.Source.Subscribe(childCtx, func(t rx.Notification) {
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
func Timeout(duration time.Duration) rx.Operator {
	return TimeoutConfigure{Duration: duration}.Use()
}
