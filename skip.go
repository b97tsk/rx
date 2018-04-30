package rx

import (
	"context"
)

type skipOperator struct {
	source Operator
	count  int
}

func (op skipOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	count := op.count
	mutable := MutableObserver{}
	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if count > 1 {
				count--
			} else {
				mutable.Observer = ob
			}
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	})
	return op.source.Call(ctx, &mutable)
}

// Skip creates an Observable that skips the first count items emitted by the
// source Observable.
func (o Observable) Skip(count int) Observable {
	if count <= 0 {
		return o
	}
	op := skipOperator{
		source: o.Op,
		count:  count,
	}
	return Observable{op}
}
