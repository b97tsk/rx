package rx

import (
	"context"
)

func never(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// Never returns an Observable that never emits anything.
func Never() Observable {
	return never
}
