package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/queue"
	"github.com/b97tsk/rx/x/schedule"
)

type observeOnObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

func (obs observeOnObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Queue queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			x.Queue.PushBack(t)
			schedule.ScheduleOnce(ctx, obs.Duration, func() {
				if x, ok := <-cx; ok {
					switch t := x.Queue.PopFront().(rx.Notification); {
					case t.HasValue:
						sink(t)
						cx <- x
					default:
						close(cx)
						sink(t)
					}
				}
			})
			cx <- x
		}
	})
}

// ObserveOn creates an Observable that emits each notification from the source
// Observable after waits for the duration to elapse.
func ObserveOn(d time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := observeOnObservable{source, d}
		return rx.Create(obs.Subscribe)
	}
}
