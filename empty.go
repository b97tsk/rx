package rx

import (
	"context"
)

func empty(ctx context.Context, sink Observer) {
	sink.Complete()
}

// Empty returns an Observable that emits no items to the Observer and
// immediately completes.
func Empty() Observable {
	return empty
}
