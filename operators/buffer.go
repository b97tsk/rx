package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type bufferObservable struct {
	Source          rx.Observable
	ClosingNotifier rx.Observable
}

func (obs bufferObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Buffers []interface{}
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.ClosingNotifier.Subscribe(ctx, func(t rx.Notification) {
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
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
}

// Buffer buffers the source Observable values until closingNotifier emits.
//
// Buffer collects values from the past as a slice, and emits that slice
// only when another Observable emits.
func Buffer(closingNotifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := bufferObservable{source, closingNotifier}
		return rx.Create(obs.Subscribe)
	}
}
