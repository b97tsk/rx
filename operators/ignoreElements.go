package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func ignoreElements(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
		return source.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}
			sink(t)
		})
	}
}

// IgnoreElements creates an Observable that ignores all values emitted by the
// source Observable and only passes ERROR or COMPLETE emission.
func IgnoreElements() rx.Operator {
	return ignoreElements
}
