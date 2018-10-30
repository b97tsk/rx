package rx

import (
	"context"
)

type takeUntilOperator struct {
	Notifier Observable
}

func (op takeUntilOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sink = Mutex(sink)

	op.Notifier.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink.Complete()
		default:
			sink(t)
		}
		cancel()
	})

	select {
	case <-ctx.Done():
		return ctx, cancel
	default:
	}

	source.Subscribe(ctx, Finally(sink, cancel))

	return ctx, cancel
}

// TakeUntil creates an Observable that emits the values emitted by the source
// Observable until a notifier Observable emits a value.
//
// TakeUntil lets values pass until a second Observable, notifier, emits
// something. Then, it completes.
func (Operators) TakeUntil(notifier Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := takeUntilOperator{notifier}
		return source.Lift(op.Call)
	}
}
