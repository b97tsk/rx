package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// IgnoreElements ignores all values emitted by the source and only passes
// errors or completions.
func IgnoreElements() rx.Operator {
	return ignoreElements
}

func ignoreElements(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		source.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}

			sink(t)
		})
	}
}
