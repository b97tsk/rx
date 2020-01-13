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

func just(one interface{}) Observable {
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		sink.Next(one)
		sink.Complete()
		return Done()
	}
}

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice(slice []interface{}) Observable {
	switch {
	case len(slice) > 1:
		return fromSliceObservable{slice}.Subscribe
	case len(slice) == 1:
		return just(slice[0])
	default:
		return Empty()
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(values ...interface{}) Observable {
	return FromSlice(values)
}
