package rx

import (
	"context"
)

type skipLastOperator struct {
	count int
}

func (op skipLastOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		buffer     = make([]interface{}, op.count)
		bufferSize = op.count
		index      int
		count      int
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if count < bufferSize {
				count++
			} else {
				ob.Next(buffer[index])
			}
			buffer[index] = t.Value
			index = (index + 1) % bufferSize
		default:
			t.Observe(ob)
		}
	})
}

// SkipLast creates an Observable that skip the last count values emitted by
// the source Observable.
func (o Observable) SkipLast(count int) Observable {
	if count <= 0 {
		return o
	}
	op := skipLastOperator{count}
	return o.Lift(op.Call)
}
