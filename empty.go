package rx

import (
	"context"
)

// Empty returns an Observable that emits no items to the Observer and
// immediately completes.
func Empty() Observable {
	return empty
}

func empty(ctx context.Context, sink Observer) {
	sink.Complete()
}
