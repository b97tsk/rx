package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type takeLastOperator struct {
	Count int
}

func (op takeLastOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var queue queue.Queue
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if queue.Len() == op.Count {
				queue.PopFront()
			}
			queue.PushBack(t.Value)
		case t.HasError:
			sink(t)
		default:
			for i, j := 0, queue.Len(); i < j; i++ {
				if ctx.Err() != nil {
					return
				}
				sink.Next(queue.At(i))
			}
			sink(t)
		}
	})
}

// TakeLast creates an Observable that emits only the last count values emitted
// by the source Observable.
//
// TakeLast remembers the latest count values, then emits those only when the
// source completes.
func (Operators) TakeLast(count int) OperatorFunc {
	return func(source Observable) Observable {
		if count <= 0 {
			return Empty()
		}
		op := takeLastOperator{count}
		return source.Lift(op.Call)
	}
}
