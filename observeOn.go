package rx

import (
	"container/list"
	"context"
	"time"
)

type observeOnOperator struct {
	source    Operator
	delay     time.Duration
	scheduler Scheduler
}

func (op observeOnOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	try := cancellableLocker{}
	queue := list.List{}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			queue.PushBack(t)

			op.scheduler.ScheduleOnce(ctx, op.delay, func() {
				if try.Lock() {
					t := queue.Remove(queue.Front()).(Notification)
					switch {
					case t.HasValue:
						ob.Next(t.Value)
						try.Unlock()
					case t.HasError:
						try.Cancel()
						ob.Error(t.Value.(error))
						cancel()
					default:
						try.Cancel()
						ob.Complete()
						cancel()
					}
				}
			})

			try.Unlock()
		}
	}))

	return ctx, cancel
}

// ObserveOn creates an Observable that re-emits all notifications from source
// Observable with specified scheduler.
func (o Observable) ObserveOn(s Scheduler, delay time.Duration) Observable {
	op := observeOnOperator{
		source:    o.Op,
		delay:     delay,
		scheduler: s,
	}
	return Observable{op}
}
