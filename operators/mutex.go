package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func mutex(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
		return source.Subscribe(ctx, rx.Mutex(sink))
	}
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func Mutex() rx.Operator {
	return mutex
}
