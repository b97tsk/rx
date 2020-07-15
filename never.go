package rx

import (
	"context"
)

func never(context.Context, Observer) {}

// Never returns an Observable that never emits anything.
func Never() Observable {
	return never
}
