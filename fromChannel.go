package rx

import (
	"context"
)

type fromChannelObservable struct {
	Chan <-chan interface{}
}

func (obs fromChannelObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)
	for {
		select {
		case <-ctx.Done():
			return ctx, ctx.Cancel
		case val, ok := <-obs.Chan:
			if !ok {
				sink.Complete()
				ctx.Unsubscribe(Complete)
				return ctx, ctx.Cancel
			}
			sink.Next(val)
			// Check if ctx is canceled before next loop, such that
			// Take(1) would exactly take one from the channel.
			if ctx.Err() != nil {
				return ctx, ctx.Cancel
			}
		}
	}
}

// FromChannel creates an Observable that emits values from a channel, and
// completes when the channel closes.
func FromChannel(c <-chan interface{}) Observable {
	return fromChannelObservable{c}.Subscribe
}
