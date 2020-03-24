package rx

import (
	"context"
)

type fromChannelObservable struct {
	Chan <-chan interface{}
}

func (obs fromChannelObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return Done(ctx)
		case val, ok := <-obs.Chan:
			if !ok {
				sink.Complete()
				return Done(ctx)
			}
			sink.Next(val)
			// Check if ctx is canceled before next loop, such that
			// Take(1) would exactly take one from the channel.
			if ctx.Err() != nil {
				return Done(ctx)
			}
		}
	}
}

// FromChannel creates an Observable that emits values from a channel, and
// completes when the channel closes.
func FromChannel(c <-chan interface{}) Observable {
	return fromChannelObservable{c}.Subscribe
}
