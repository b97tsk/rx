package rx

import (
	"context"
)

type takeOperator struct {
	count int
}

func (op takeOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		count           = op.count
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			if count > 0 {
				count--
				if count > 0 {
					t.Observe(ob)
				} else {
					mutableObserver = NopObserver
					t.Observe(ob)
					ob.Complete()
					cancel()
				}
			} else {
				mutableObserver = NopObserver
				ob.Complete()
				cancel()
			}
		default:
			t.Observe(ob)
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Take creates an Observable that emits only the first count values emitted
// by the source Observable.
//
// Take takes the first count values from the source, then completes.
func (o Observable) Take(count int) Observable {
	op := takeOperator{count}
	return o.Lift(op.Call)
}
