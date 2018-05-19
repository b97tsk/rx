package rx

import (
	"context"
	"time"
)

type rangeOperator struct {
	low, high int
	delay     time.Duration
	scheduler Scheduler
}

func (op rangeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	if op.scheduler != nil {
		ctx, cancel := context.WithCancel(ctx)
		index := op.low

		op.scheduler.Schedule(ctx, op.delay, func() {
			if index < op.high {
				sink.Next(index)
				index++
				return
			}
			sink.Complete()
			cancel()
		})

		return ctx, cancel
	}

	done := ctx.Done()

	for index := op.low; index < op.high; index++ {
		select {
		case <-done:
			return canceledCtx, doNothing
		default:
		}
		sink.Next(index)
	}

	sink.Complete()
	return canceledCtx, doNothing
}

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	op := rangeOperator{low: low, high: high}
	return Observable{}.Lift(op.Call)
}

// RangeOn creates an Observable that emits a sequence of integers within a
// specified range, on the specified Scheduler.
func RangeOn(low, high int, s Scheduler, delay time.Duration) Observable {
	op := rangeOperator{low, high, delay, s}
	return Observable{}.Lift(op.Call)
}
