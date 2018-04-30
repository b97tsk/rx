package rx

import (
	"context"
	"time"
)

type rangeOperator struct {
	start     int
	count     int
	delay     time.Duration
	scheduler Scheduler
}

func (op rangeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	if op.scheduler != nil {
		ctx, cancel := context.WithCancel(ctx)
		index := 0

		op.scheduler.Schedule(ctx, op.delay, func() {
			if index < op.count {
				ob.Next(op.start + index)
				index++
				return
			}
			ob.Complete()
			cancel()
		})

		return ctx, cancel
	}

	done := ctx.Done()

	for index := 0; index < op.count; index++ {
		select {
		case <-done:
			return canceledCtx, noopFunc
		default:
		}
		ob.Next(op.start + index)
	}

	ob.Complete()
	return canceledCtx, noopFunc
}

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(start, count int) Observable {
	op := rangeOperator{start: start, count: count}
	return Observable{op}
}

// RangeOn creates an Observable that emits a sequence of integers within a
// specified range, on the specified Scheduler.
func RangeOn(start, count int, s Scheduler, delay time.Duration) Observable {
	op := rangeOperator{
		start:     start,
		count:     count,
		delay:     delay,
		scheduler: s,
	}
	return Observable{op}
}
