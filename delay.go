package rx

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type delayOperator struct {
	source    Operator
	timeout   time.Duration
	scheduler Scheduler
}

type delayValue struct {
	Time time.Time
	Notification
}

func (op delayOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case schedulerOption:
			op.scheduler = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op delayOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	mu := sync.Mutex{}
	queue := list.List{}
	scheduleCtx := canceledCtx
	scheduleDone := scheduleCtx.Done()

	var doSchedule func(time.Duration)

	doSchedule = func(timeout time.Duration) {
		select {
		case <-scheduleDone:
		default:
			return
		}

		scheduleCtx, _ = op.scheduler.ScheduleOnce(ctx, timeout, func() {
			mu.Lock()
			defer mu.Unlock()

			for e := queue.Front(); e != nil; e, _ = e.Next(), queue.Remove(e) {
				select {
				case <-done:
					return
				default:
				}
				t := e.Value.(delayValue)
				now := op.scheduler.Now()
				if t.Time.After(now) {
					doSchedule(t.Time.Sub(now))
					return
				}
				switch {
				case t.HasValue:
					ob.Next(t.Value)
				case t.HasError:
					ob.Error(t.Value.(error))
					cancel()
				default:
					ob.Complete()
					cancel()
				}
			}
		})
		scheduleDone = scheduleCtx.Done()
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case t.HasValue:
			queue.PushBack(delayValue{
				Time:         op.scheduler.Now().Add(op.timeout),
				Notification: t,
			})
			doSchedule(op.timeout)
		case t.HasError:
			// Error notification will not be delayed.
			queue.Init()
			ob.Error(t.Value.(error))
			cancel()
		default:
			queue.PushBack(delayValue{
				Time: op.scheduler.Now().Add(op.timeout),
			})
			doSchedule(op.timeout)
		}
	}))

	return ctx, cancel
}

// Delay delays the emission of items from the source Observable by a given
// timeout.
func (o Observable) Delay(timeout time.Duration) Observable {
	op := delayOperator{
		source:    o.Op,
		timeout:   timeout,
		scheduler: DefaultScheduler,
	}
	return Observable{op}
}
