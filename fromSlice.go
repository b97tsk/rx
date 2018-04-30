package rx

import (
	"context"
	"time"
)

type fromSliceOperator struct {
	slice     []interface{}
	delay     time.Duration
	scheduler Scheduler
}

func (op fromSliceOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	if op.scheduler != nil {
		ctx, cancel := context.WithCancel(ctx)
		index := 0

		op.scheduler.Schedule(ctx, op.delay, func() {
			if index < len(op.slice) {
				val := op.slice[index]
				ob.Next(val)
				index++
				return
			}
			ob.Complete()
			cancel()
		})

		return ctx, cancel
	}

	done := ctx.Done()

	for _, val := range op.slice {
		select {
		case <-done:
			return canceledCtx, noopFunc
		default:
		}
		ob.Next(val)
	}
	ob.Complete()

	return canceledCtx, noopFunc
}

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice(slice []interface{}) Observable {
	if len(slice) == 0 {
		return Empty()
	}
	op := fromSliceOperator{slice: slice}
	return Observable{op}
}

// FromSliceOn creates an Observable that emits values from a slice, one after
// the other, and then completes, on the specified Scheduler.
func FromSliceOn(slice []interface{}, s Scheduler, delay time.Duration) Observable {
	if len(slice) == 0 {
		return EmptyOn(s, delay)
	}
	op := fromSliceOperator{
		slice:     slice,
		delay:     delay,
		scheduler: s,
	}
	return Observable{op}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(slice ...interface{}) Observable {
	if len(slice) == 0 {
		return Empty()
	}
	op := fromSliceOperator{slice: slice}
	return Observable{op}
}
