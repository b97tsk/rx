package rx

import (
	"context"
)

type skipLastOperator struct {
	Count int
}

func (op skipLastOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		buffer     = make([]interface{}, op.Count)
		bufferSize = op.Count
		index      int
		count      int
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if count < bufferSize {
				count++
			} else {
				sink.Next(buffer[index])
			}
			buffer[index] = t.Value
			index = (index + 1) % bufferSize
		default:
			sink(t)
		}
	})
}

// SkipLast creates an Observable that skip the last count values emitted by
// the source Observable.
func (Operators) SkipLast(count int) OperatorFunc {
	return func(source Observable) Observable {
		if count <= 0 {
			return source
		}
		op := skipLastOperator{count}
		return source.Lift(op.Call)
	}
}
