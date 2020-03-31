package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/schedule"
)

type debounceTimeObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

func (obs debounceTimeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var scheduleCancel context.CancelFunc

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				x.HasLatestValue = true

				cx <- x

				if scheduleCancel != nil {
					scheduleCancel()
				}

				_, scheduleCancel = schedule.ScheduleOnce(ctx, obs.Duration, func() {
					if x, ok := <-cx; ok {
						if x.HasLatestValue {
							sink.Next(x.LatestValue)
							x.HasLatestValue = false
						}
						cx <- x
					}
				})

			case t.HasError:
				close(cx)
				sink(t)

			default:
				close(cx)
				if x.HasLatestValue {
					sink.Next(x.LatestValue)
				}
				sink(t)
			}
		}
	})
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
//
// It's like Delay, but passes only the most recent value from each burst of
// emissions.
func DebounceTime(duration time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := debounceTimeObservable{source, duration}
		return rx.Create(obs.Subscribe)
	}
}
