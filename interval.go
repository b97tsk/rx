package rx

import (
	"context"
	"time"
)

type intervalOperator struct {
	InitialDelay time.Duration
	Period       time.Duration
}

func (op intervalOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	if op.InitialDelay != op.Period {
		scheduleOnce(ctx, op.InitialDelay, func() {
			index := 0
			wait := make(chan struct{})
			schedule(ctx, op.Period, func() {
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
		schedule(ctx, op.Period, func() {
			sink.Next(index)
			index++
		})
	}

	return ctx, cancel
}

// Interval creates an Observable that emits sequential integers every
// specified interval of time.
func Interval(period time.Duration) Observable {
	op := intervalOperator{period, period}
	return Empty().Lift(op.Call)
}

// Timer creates an Observable that starts emitting after an initialDelay and
// emits ever increasing integers after each period of time thereafter.
//
// Its like Interval, but you can specify when should the emissions start.
func Timer(initialDelay, period time.Duration) Observable {
	op := intervalOperator{initialDelay, period}
	return Empty().Lift(op.Call)
}
