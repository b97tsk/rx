package rx

import (
	"context"
)

// Go creates an Observable that mirrors the source Observable in a goroutine.
func (Operators) Go() Operator {
	return func(source Observable) Observable {
		return Create(
			func(ctx context.Context, sink Observer) {
				go source.Subscribe(ctx, sink)
			},
		)
	}
}
