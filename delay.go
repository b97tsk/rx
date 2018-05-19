package rx

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type delayOperator struct {
	timeout time.Duration
}

type delayValue struct {
	Time time.Time
	Notification
}

func (op delayOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	scheduleCtx := canceledCtx
	scheduleDone := scheduleCtx.Done()

	var (
		mu         sync.Mutex
		queue      list.List
		doSchedule func(time.Duration)
	)

	doSchedule = func(timeout time.Duration) {
		select {
		case <-scheduleDone:
		default:
			return
		}

		scheduleCtx, _ = scheduleOnce(ctx, timeout, func() {
			mu.Lock()
			defer mu.Unlock()

			for e := queue.Front(); e != nil; e, _ = e.Next(), queue.Remove(e) {
				select {
				case <-done:
					return
				default:
				}
				t := e.Value.(delayValue)
				now := time.Now()
				if t.Time.After(now) {
					doSchedule(t.Time.Sub(now))
					return
				}
				switch {
				case t.HasValue:
					sink(t.Notification)
				default:
					sink(t.Notification)
					cancel()
				}
			}
		})
		scheduleDone = scheduleCtx.Done()
	}

	source.Subscribe(ctx, func(t Notification) {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case t.HasValue:
			queue.PushBack(delayValue{
				Time:         time.Now().Add(op.timeout),
				Notification: t,
			})
			doSchedule(op.timeout)
		case t.HasError:
			// Error notification will not be delayed.
			queue.Init()
			sink(t)
			cancel()
		default:
			queue.PushBack(delayValue{
				Time: time.Now().Add(op.timeout),
			})
			doSchedule(op.timeout)
		}
	})

	return ctx, cancel
}

// Delay delays the emission of items from the source Observable by a given
// timeout.
func (o Observable) Delay(timeout time.Duration) Observable {
	op := delayOperator{timeout}
	return o.Lift(op.Call)
}
