package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/x/schedule"
)

type intervalObservable struct {
	InitialDelay time.Duration
	Period       time.Duration
}

func (obs intervalObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	if obs.InitialDelay != obs.Period {
		schedule.Once(ctx, obs.InitialDelay, func() {
			index := 0
			wait := make(chan struct{})
			schedule.Periodic(ctx, obs.Period, func() {
				<-wait
				sink.Next(index)
				index++
			})
			sink.Next(index)
			index++
			close(wait)
		})
	} else {
		index := 0
		schedule.Periodic(ctx, obs.Period, func() {
			sink.Next(index)
			index++
		})
	}

	return ctx, cancel
}

// Interval creates an Observable that emits sequential integers every
// specified interval of time.
func Interval(period time.Duration) Observable {
	return intervalObservable{period, period}.Subscribe
}

// Timer creates an Observable that starts emitting after an initialDelay and
// emits ever increasing integers after each period of time thereafter.
//
// Its like Interval, but you can specify when should the emissions start.
func Timer(initialDelay, period time.Duration) Observable {
	return intervalObservable{initialDelay, period}.Subscribe
}
