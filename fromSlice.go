package rx

import (
	"context"
)

type fromSliceObservable struct {
	Slice []interface{}
}

func (obs fromSliceObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	for _, val := range obs.Slice {
		if ctx.Err() != nil {
			return Done()
		}
		sink.Next(val)
	}
	sink.Complete()
	return Done()
}

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice(slice []interface{}) Observable {
	if len(slice) == 0 {
		return Empty()
	}
	return fromSliceObservable{slice}.Subscribe
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(values ...interface{}) Observable {
	return FromSlice(values)
}
