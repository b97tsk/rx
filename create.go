package rx

import (
	"context"
)

// Create creates a new Observable, that will execute the specified function
// when an Observer subscribes to it.
//
// Create custom Observable, that does whatever you like.
func Create(subscribe func(context.Context, Observer) (context.Context, context.CancelFunc)) Observable {
	return Observable{}.Lift(
		func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
			return subscribe(ctx, sink)
		},
	)
}
