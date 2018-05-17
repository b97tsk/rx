package rx

import (
	"container/list"
	"context"
	"time"
)

type observeOnOperator struct {
	delay     time.Duration
	scheduler Scheduler
}

func (op observeOnOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		queue list.List
		try   cancellableLocker
	)

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			defer try.Unlock()
			queue.PushBack(t)
			op.scheduler.ScheduleOnce(ctx, op.delay, func() {
				if try.Lock() {
					t := queue.Remove(queue.Front()).(Notification)
					switch {
					case t.HasValue:
						defer try.Unlock()
						ob.Next(t.Value)
					case t.HasError:
						try.CancelAndUnlock()
						ob.Error(t.Value.(error))
						cancel()
					default:
						try.CancelAndUnlock()
						ob.Complete()
						cancel()
					}
				}
			})
		}
	})

	return ctx, cancel
}

// ObserveOn creates an Observable that re-emits all notifications from source
// Observable with specified scheduler.
func (o Observable) ObserveOn(s Scheduler, delay time.Duration) Observable {
	op := observeOnOperator{delay, s}
	return o.Lift(op.Call)
}
