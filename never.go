package rx

import (
	"context"
)

// Never creates an Observable that never emits anything.
func Never() Observable {
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		return context.WithCancel(ctx)
	}
}
