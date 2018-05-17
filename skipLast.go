package rx

import (
	"context"
)

type skipLastOperator struct {
	source Operator
	count  int
}

func (op skipLastOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	buffer := make([]interface{}, op.count)
	bufferSize := op.count
	index := 0
	count := 0
	return op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if count < bufferSize {
				count++
			} else {
				ob.Next(buffer[index])
			}
			buffer[index] = t.Value
			index = (index + 1) % bufferSize
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	})
}

// SkipLast creates an Observable that skip the last count values emitted by
// the source Observable.
func (o Observable) SkipLast(count int) Observable {
	if count <= 0 {
		return o
	}
	op := skipLastOperator{
		source: o.Op,
		count:  count,
	}
	return Observable{op}
}
