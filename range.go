package rx

import (
	"context"
)

type rangeObservable struct {
	Low, High int
}

func (obs rangeObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	for index := obs.Low; index < obs.High; index++ {
		if ctx.Err() != nil {
			return Done(ctx)
		}
		sink.Next(index)
	}
	sink.Complete()
	return Done(ctx)
}

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	if low >= high {
		return Empty()
	}
	return rangeObservable{low, high}.Subscribe
}
