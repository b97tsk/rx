package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

// Timeout creates an Observable that mirrors the source Observable or throws
// rx.ErrTimeout if the source does not emit a value in given time span.
func Timeout(d time.Duration) rx.Operator {
	return TimeoutConfigure{Duration: d}.Make()
}

// A TimeoutConfigure is a configure for Timeout.
type TimeoutConfigure struct {
	Duration   time.Duration
	Observable rx.Observable
}

// Make creates an Operator from this configure.
func (configure TimeoutConfigure) Make() rx.Operator {
	if configure.Observable == nil {
		configure.Observable = rx.Throw(rx.ErrTimeout)
	}
	return func(source rx.Observable) rx.Observable {
		return timeoutObservable{source, configure}.Subscribe
	}
}

type timeoutObservable struct {
	Source rx.Observable
	TimeoutConfigure
}

func (obs timeoutObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	childCtx, childCancel := context.WithCancel(ctx)

	race := make(chan struct{}, 1)
	race <- struct{}{}

	var scheduleCancel context.CancelFunc

	obsTimer := rx.Timer(obs.Duration)
	doSchedule := func() {
		if scheduleCancel != nil {
			scheduleCancel()
		}
		var scheduleCtx context.Context
		scheduleCtx, scheduleCancel = context.WithCancel(childCtx)
		obsTimer.Subscribe(scheduleCtx, func(t rx.Notification) {
			if t.HasValue {
				return
			}
			if _, ok := <-race; ok {
				close(race)
				childCancel()
				obs.Observable.Subscribe(ctx, sink)
			}
		})
	}

	doSchedule()

	obs.Source.Subscribe(childCtx, func(t rx.Notification) {
		if x, ok := <-race; ok {
			switch {
			case t.HasValue:
				sink(t)
				race <- x
				doSchedule()
			default:
				close(race)
				childCancel()
				sink(t)
			}
		}
	})
}
