package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Go creates an Observable that mirrors the source Observable in a goroutine.
func Go() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return rx.Create(
			func(ctx context.Context, sink rx.Observer) {
				go source.Subscribe(ctx, sink)
			},
		)
	}
}
