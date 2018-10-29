package rx

import (
	"context"
)

type retryWhenOperator struct {
	Notifier func(Observable) Observable
}

func (op retryWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sourceCtx, sourceCancel := canceledCtx, nothingToDo

	var (
		subject  *Subject
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			if subject == nil {
				subject = NewSubject()
				obsv := op.Notifier(subject.Observable)
				obsv.Subscribe(ctx, func(t Notification) {
					switch {
					case t.HasValue:
						sourceCancel()

						sourceCtx, sourceCancel = context.WithCancel(ctx)
						source.Subscribe(sourceCtx, observer)

					default:
						sink(t)
						cancel()
					}
				})
			}
			subject.Next(t.Value.(error))
		default:
			sink(t)
			cancel()
		}
	}

	sourceCtx, sourceCancel = context.WithCancel(ctx)
	source.Subscribe(sourceCtx, observer)

	return ctx, cancel
}

// RetryWhen creates an Observable that mirrors the source Observable with the
// exception of an Error. If the source Observable calls Error, this method
// will emit the Throwable that caused the Error to the Observable returned
// from notifier. If that Observable calls Complete or Error then this method
// will call Complete or Error on the child subscription. Otherwise this method
// will resubscribe to the source Observable.
func (Operators) RetryWhen(notifier func(Observable) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := retryWhenOperator{notifier}
		return source.Pipe(MakeFunc(op.Call), operators.Mutex())
	}
}
