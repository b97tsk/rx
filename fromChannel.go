package rx

import (
	"context"
)

type fromChannelOperator struct {
	Chan <-chan interface{}
}

func (op fromChannelOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	done := ctx.Done()
	for {
		select {
		case <-done:
			return canceledCtx, nothingToDo
		case val, ok := <-op.Chan:
			if !ok {
				sink.Complete()
				return canceledCtx, nothingToDo
			}
			sink.Next(val)
			// Check done before next loop, such that Take(1)
			// would exactly take one from the channel.
			select {
			case <-done:
				return canceledCtx, nothingToDo
			default:
			}
		}
	}
}

// FromChannel creates an Observable that emits values from a channel, and
// completes when the channel closes.
func FromChannel(c <-chan interface{}) Observable {
	op := fromChannelOperator{c}
	return Observable{}.Lift(op.Call)
}
