package rx

import (
	"context"
)

// FromChannel creates an Observable that emits values from a channel, and
// completes when the channel closes.
func FromChannel(c <-chan interface{}) Observable {
	return Create(
		func(ctx context.Context, sink Observer) {
			done := ctx.Done()
			for {
				select {
				case <-done:
					return
				case val, ok := <-c:
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
		},
	)
}
