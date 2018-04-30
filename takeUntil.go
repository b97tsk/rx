package rx

import (
	"context"
)

type takeUntilOperator struct {
	source   Operator
	notifier Observable
}

func (op takeUntilOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	ob = Normalize(ob)

	op.notifier.Subscribe(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			ob.Complete()
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
		cancel()
	}))

	select {
	case <-done:
		return ctx, cancel
	default:
	}

	op.source.Call(ctx, withFinalizer(ob, cancel))

	return ctx, cancel
}

// TakeUntil creates an Observable that emits the values emitted by the source
// Observable until a notifier Observable emits a value.
//
// TakeUntil lets values pass until a second Observable, notifier, emits
// something. Then, it completes.
func (o Observable) TakeUntil(notifier Observable) Observable {
	op := takeUntilOperator{
		source:   o.Op,
		notifier: notifier,
	}
	return Observable{op}
}
