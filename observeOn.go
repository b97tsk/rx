package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/x/queue"
)

type observeOnObservable struct {
	Source   Observable
	Duration time.Duration
}

func (obs observeOnObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Queue queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			x.Queue.PushBack(t)
			scheduleOnce(ctx, obs.Duration, func() {
				if x, ok := <-cx; ok {
					switch t := x.Queue.PopFront().(Notification); {
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

	return ctx, cancel
}

// ObserveOn creates an Observable that emits each notification from the source
// Observable after waits for the duration to elapse.
func (Operators) ObserveOn(d time.Duration) Operator {
	return func(source Observable) Observable {
		return observeOnObservable{source, d}.Subscribe
	}
}
