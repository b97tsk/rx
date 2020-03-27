package rx

import (
	"context"
)

type fromChannelObservable struct {
	Chan <-chan interface{}
}

func (obs fromChannelObservable) Subscribe(ctx context.Context, sink Observer) {
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-obs.Chan:
			if !ok {
				sink.Complete()
				return
			}
			sink.Next(val)
			// Check if ctx is canceled before next loop, such that
			// Take(1) would exactly take one from the channel.
			if ctx.Err() != nil {
				return
			}
		}
	}
}

// FromChannel creates an Observable that emits values from a channel, and
// completes when the channel closes.
func FromChannel(c <-chan interface{}) Observable {
	obs := fromChannelObservable{c}
	return Create(obs.Subscribe)
}
