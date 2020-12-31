package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Go mirrors the source in a goroutine.
func Go() rx.Operator {
	return go1
}

func go1(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		go source.Subscribe(ctx, sink)
	}
}
