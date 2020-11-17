package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Pairwise creates an Observable that groups pairs of consecutive emissions
// together and emits them as rx.Pairs.
func Pairwise() rx.Operator {
	return pairwise
}

func pairwise(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		var p struct {
			Value    interface{}
			HasValue bool
		}

		source.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				if p.HasValue {
					sink.Next(rx.Pair{
						First:  p.Value,
						Second: t.Value,
					})
					p.Value = t.Value
				} else {
					p.Value = t.Value
					p.HasValue = true
				}
			default:
				sink(t)
			}
		})
	}
}
