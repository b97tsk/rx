package rx

import (
	"context"
)

// Never returns an Observable that never emits anything.
func Never() Observable {
	return never
}

func never(context.Context, Observer) {}
