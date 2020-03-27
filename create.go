package rx

import (
	"context"
)

// Create creates a new Observable, that will execute the specified function
// when an Observer subscribes to it.
//
// It's the caller's responsibility to follow the rule of Observable that
// no more emissions pass to the sink Observer after an ERROR or COMPLETE
// emission passes to it.
func Create(subscribe func(context.Context, Observer)) Observable {
	return func(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx := NewContext(parent)
		sink = DoAtLast(sink, ctx.AtLast)
		subscribe(ctx, sink)
		return ctx, ctx.Cancel
	}
}
