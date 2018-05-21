package rx

import (
	"context"
)

type rangeOperator struct {
	Low, High int
}

func (op rangeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	done := ctx.Done()
	for index := op.Low; index < op.High; index++ {
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
	op := rangeOperator{low, high}
	return Observable{}.Lift(op.Call)
}
