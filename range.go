package rx

import (
	"context"
)

type rangeObservable struct {
	Low, High int
}

func (obs rangeObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)
	for index := obs.Low; index < obs.High; index++ {
		if ctx.Err() != nil {
			return ctx, ctx.Cancel
		}
		sink.Next(index)
	}
	sink.Complete()
	ctx.Unsubscribe(Complete)
	return ctx, ctx.Cancel
}

// Range creates an Observable that emits a sequence of integers within a
// specified range.
func Range(low, high int) Observable {
	if low >= high {
		return Empty()
	}
	return rangeObservable{low, high}.Subscribe
}
