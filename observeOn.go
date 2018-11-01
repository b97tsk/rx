package rx

import (
	"container/list"
	"context"
	"time"
)

type observeOnOperator struct {
	Duration time.Duration
}

func (op observeOnOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		queue list.List
		try   cancellableLocker
	)

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			defer try.Unlock()
			queue.PushBack(t)
			scheduleOnce(ctx, op.Duration, func() {
				if try.Lock() {
					switch t := queue.Remove(queue.Front()).(Notification); {
					case t.HasValue:
						defer try.Unlock()
						sink(t)
					default:
						try.CancelAndUnlock()
						sink(t)
					}
				}
			})
		}
	})

	return ctx, cancel
}

// ObserveOn creates an Observable that emits each notification from the source
// Observable after waits for the duration to elapse.
func (Operators) ObserveOn(d time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := observeOnOperator{d}
		return source.Lift(op.Call)
	}
}
