package rx

import (
	"container/list"
	"context"
	"sync"
	"time"
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

	var (
		scheduleCtx  = canceledCtx

		mutex      sync.Mutex
		queue      list.List
		doSchedule func(time.Duration)
	)

	doSchedule = func(timeout time.Duration) {
		if !isDone(scheduleCtx) {
			return
		}

		scheduleCtx, _ = scheduleOnce(ctx, timeout, func() {
			mutex.Lock()
			defer mutex.Unlock()
			for e := queue.Front(); e != nil; e, _ = e.Next(), queue.Remove(e) {
				if isDone(ctx) {
					return
				}
				t := e.Value.(delayValue)
				now := time.Now()
				if t.Time.After(now) {
					doSchedule(t.Time.Sub(now))
					return
				}
				sink(t.Notification)
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		mutex.Lock()
		switch {
		case t.HasValue:
			queue.PushBack(delayValue{
				Time:         time.Now().Add(op.Duration),
				Notification: t,
			})
			doSchedule(op.Duration)
		case t.HasError:
			// Error notification will not be delayed.
			queue.Init()
			sink(t)
		default:
			queue.PushBack(delayValue{
				Time: time.Now().Add(op.Duration),
			})
			doSchedule(op.Duration)
		}
		mutex.Unlock()
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
