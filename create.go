package rx

import (
	"context"
)

// Create creates a new Observable, that will execute the specified function
// when an Observer subscribes to it.
//
// It's the caller's responsibility to ensure no more emissions passed to
// the sink Observer after an ERROR or COMPLETE emission passes to it.
func Create(subscribe func(context.Context, Observer) (context.Context, context.CancelFunc)) Observable {
	return Observable{}.Lift(
		func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
			return subscribe(ctx, sink)
		},
	)
}
