package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Timeout creates an Observable that mirrors the source Observable or throws
// rx.ErrTimeout if the source does not emit a value in given time span.
func Timeout(d time.Duration) rx.Operator {
	return TimeoutWith(d, rx.Throw(rx.ErrTimeout))
}

// TimeoutWith creates an Observable that mirrors the source Observable or
// specified Observable if the source does not emit a value in given time span.
func TimeoutWith(d time.Duration, switchTo rx.Observable) rx.Operator {
	if switchTo == nil {
		panic("TimeoutWith: switchTo is nil")
	}

	return func(source rx.Observable) rx.Observable {
		return timeoutObservable{source, d, switchTo}.Subscribe
	}
}

type timeoutObservable struct {
	Source   rx.Observable
	Duration time.Duration
	SwitchTo rx.Observable
}

func (obs timeoutObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	childCtx, childCancel := context.WithCancel(ctx)

	var race critical.Section

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

			if critical.Enter(&race) {
				critical.Close(&race)
				childCancel()
				obs.SwitchTo.Subscribe(ctx, sink)
			}
		})
	}

	doSchedule()

	obs.Source.Subscribe(childCtx, func(t rx.Notification) {
		if critical.Enter(&race) {
			switch {
			case t.HasValue:
				sink(t)
				critical.Leave(&race)
				doSchedule()
			default:
				critical.Close(&race)
				childCancel()
				sink(t)
			}
		}
	})
}
