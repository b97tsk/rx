package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/x/queue"
)

type delayOperator struct {
	Duration time.Duration
}

type delayValue struct {
	Time time.Time
	Notification
}

func (op delayOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

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

	source.Subscribe(ctx, func(t Notification) {
		x := <-cx
		switch {
		case t.HasValue:
			x.Queue.PushBack(
				delayValue{
					Time:         time.Now().Add(op.Duration),
					Notification: t,
				},
			)
			if !x.Scheduled {
				x.Scheduled = true
				doSchedule(op.Duration)
			}
		case t.HasError:
			// ERROR notification will not be delayed.
			x.Queue.Init()
			sink(t)
		default:
			x.Queue.PushBack(
				delayValue{
					Time: time.Now().Add(op.Duration),
				},
			)
			if !x.Scheduled {
				x.Scheduled = true
				doSchedule(op.Duration)
			}
		}
		cx <- x
	})

	return ctx, cancel
}

// Delay delays the emission of items from the source Observable by a given
// timeout.
func (Operators) Delay(timeout time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := delayOperator{timeout}
		return source.Lift(op.Call)
	}
}
