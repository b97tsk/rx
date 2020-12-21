package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Buffer buffers the source Observable values until closingNotifier emits.
//
// Buffer collects values from the past as a slice, and emits that slice
// only when another Observable emits.
func Buffer(closingNotifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return bufferObservable{source, closingNotifier}.Subscribe
	}
}

type bufferObservable struct {
	Source          rx.Observable
	ClosingNotifier rx.Observable
}

func (obs bufferObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Buffer []interface{}
	}

	obs.ClosingNotifier.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				defer critical.Leave(&x.Section)

				sink.Next(x.Buffer)

				x.Buffer = nil

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				sink(t)
			}
		}
	})

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Buffer = append(x.Buffer, t.Value)

				critical.Leave(&x.Section)

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				sink(t)
			}
		}
	})
}
