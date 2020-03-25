package rx

import (
	"context"
)

type bufferObservable struct {
	Source          Observable
	ClosingNotifier Observable
}

func (obs bufferObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	type X struct {
		Buffers []interface{}
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.ClosingNotifier.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sink.Next(x.Buffers)
				x.Buffers = nil
				cx <- x
			default:
				close(cx)
				sink(t)
			}
		}
	})

	if ctx.Err() != nil {
		return ctx, ctx.Cancel
	}

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Buffers = append(x.Buffers, t.Value)
				cx <- x
			default:
				close(cx)
				sink(t)
			}
		}
	})

	return ctx, ctx.Cancel
}

// Buffer buffers the source Observable values until closingNotifier emits.
//
// Buffer collects values from the past as a slice, and emits that slice
// only when another Observable emits.
func (Operators) Buffer(closingNotifier Observable) Operator {
	return func(source Observable) Observable {
		return bufferObservable{source, closingNotifier}.Subscribe
	}
}
