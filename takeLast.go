package rx

import (
	"container/list"
	"context"
)

type takeLastOperator struct {
	Count int
}

func (op takeLastOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var buffer list.List
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if buffer.Len() >= op.Count {
				buffer.Remove(buffer.Front())
			}
			buffer.PushBack(t.Value)
		case t.HasError:
			sink(t)
		default:
			for e := buffer.Front(); e != nil; e = e.Next() {
				sink.Next(e.Value)
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
