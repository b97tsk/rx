package rx

import (
	"context"
)

type skipOperator struct {
	count int
}

func (op skipOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		count           = op.count
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			if count > 1 {
				count--
			} else {
				mutableObserver = ob
			}
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	}

	return source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })
}

// Skip creates an Observable that skips the first count items emitted by the
// source Observable.
func (o Observable) Skip(count int) Observable {
	if count <= 0 {
		return o
	}
	op := skipOperator{count}
	return o.Lift(op.Call)
}
