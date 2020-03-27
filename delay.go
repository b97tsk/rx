package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/x/queue"
)

type delayObservable struct {
	Source   Observable
	Duration time.Duration
}

type delayValue struct {
	Time time.Time
	Notification
}

func (obs delayObservable) Subscribe(ctx context.Context, sink Observer) {
	type X struct {
		Queue     queue.Queue
		Scheduled bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doSchedule func(time.Duration)

	doSchedule = func(timeout time.Duration) {
		scheduleOnce(ctx, timeout, func() {
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

	obs.Source.Subscribe(ctx, func(t Notification) {
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
func (Operators) Delay(timeout time.Duration) Operator {
	return func(source Observable) Observable {
		obs := delayObservable{source, timeout}
		return Create(obs.Subscribe)
	}
}
