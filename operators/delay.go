package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/queue"
	"github.com/b97tsk/rx/x/schedule"
)

type delayObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

type delayValue struct {
	Time time.Time
	rx.Notification
}

func (obs delayObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Queue     queue.Queue
		Scheduled bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doSchedule func(time.Duration)

	doSchedule = func(timeout time.Duration) {
		schedule.ScheduleOnce(ctx, timeout, func() {
			x := <-cx
			x.Scheduled = false
			for x.Queue.Len() > 0 {
				if ctx.Err() != nil {
					break
				}
				t := x.Queue.Front().(delayValue)
				now := time.Now()
				if t.Time.After(now) {
					x.Scheduled = true
					doSchedule(t.Time.Sub(now))
					break
				}
				x.Queue.PopFront()
				sink(t.Notification)
			}
			cx <- x
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		x := <-cx
		switch {
		case t.HasValue:
			x.Queue.PushBack(
				delayValue{
					Time:         time.Now().Add(obs.Duration),
					Notification: t,
				},
			)
			if !x.Scheduled {
				x.Scheduled = true
				doSchedule(obs.Duration)
			}
		case t.HasError:
			// ERROR notification will not be delayed.
			x.Queue.Init()
			sink(t)
		default:
			x.Queue.PushBack(
				delayValue{
					Time: time.Now().Add(obs.Duration),
				},
			)
			if !x.Scheduled {
				x.Scheduled = true
				doSchedule(obs.Duration)
			}
		}
		cx <- x
	})
}

// Delay delays the emission of items from the source Observable by a given
// timeout.
func Delay(timeout time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := delayObservable{source, timeout}
		return rx.Create(obs.Subscribe)
	}
}
