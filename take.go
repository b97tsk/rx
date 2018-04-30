package rx

import (
	"context"
)

type takeOperator struct {
	source Operator
	count  int
}

func (op takeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	count := op.count

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if count > 0 {
				count--
				if count > 0 {
					ob.Next(t.Value)
				} else {
					mutable.Observer = NopObserver
					ob.Next(t.Value)
					ob.Complete()
					cancel()
				}
			} else {
				mutable.Observer = NopObserver
				ob.Complete()
				cancel()
			}
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			ob.Complete()
			cancel()
		}
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// Take creates an Observable that emits only the first count values emitted
// by the source Observable.
//
// Take takes the first count values from the source, then completes.
func (o Observable) Take(count int) Observable {
	op := takeOperator{
		source: o.Op,
		count:  count,
	}
	return Observable{op}
}
