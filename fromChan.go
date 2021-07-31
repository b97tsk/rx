package rx

import (
	"context"
)

// FromChan creates an Observable that emits values received from a channel,
// and completes when the channel closes.
func FromChan(c <-chan interface{}) Observable {
	return func(ctx context.Context, sink Observer) {
		done := ctx.Done()

		for ctx.Err() == nil {
			select {
			case <-done:
				return

			case val, ok := <-c:
				if !ok {
					sink.Complete()
					return
				}

				sink.Next(val)
			}
		}
	}
}
