package rx

import (
	"context"
)

type rangeOperator struct {
	Low, High int
}

func (op rangeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	for index := op.Low; index < op.High; index++ {
		if ctx.Err() != nil {
			return Done()
		}
		sink.Next(index)
	}
	sink.Complete()
	return Done()
}

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	switch n := high - low; {
	case n > 1:
		op := rangeOperator{low, high}
		return Observable{}.Lift(op.Call)
	case n == 1:
		return just(low)
	default:
		return Empty()
	}
}
