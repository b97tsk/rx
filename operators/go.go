package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func go1(source rx.Observable) rx.Observable {
	return rx.Create(
		func(ctx context.Context, sink rx.Observer) {
			go source.Subscribe(ctx, sink)
		},
	)
}

// Go creates an Observable that mirrors the source Observable in a goroutine.
func Go() rx.Operator {
	return go1
}
