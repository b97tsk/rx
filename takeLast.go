package rx

import (
	"container/list"
	"context"
)

type takeLastOperator struct {
	source Operator
	count  int
}

func (op takeLastOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	buffer := list.List{}
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if buffer.Len() >= op.count {
				buffer.Remove(buffer.Front())
			}
			buffer.PushBack(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			for e := buffer.Front(); e != nil; e = e.Next() {
				ob.Next(e.Value)
			}
			ob.Complete()
		}
	}))
}

// TakeLast creates an Observable that emits only the last count values emitted
// by the source Observable.
//
// TakeLast remembers the latest count values, then emits those only when the
// source completes.
func (o Observable) TakeLast(count int) Observable {
	if count <= 0 {
		return Empty()
	}
	op := takeLastOperator{
		source: o.Op,
		count:  count,
	}
	return Observable{op}
}
