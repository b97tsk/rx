package rx

import (
	"context"
)

type repeatWhenOperator struct {
	notifier func(Observable) Observable
}

func (op repeatWhenOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sourceCtx, sourceCancel := canceledCtx, noopFunc

	var (
		subject  *Subject
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			t.Observe(ob)
		case t.HasError:
			t.Observe(ob)
			cancel()
		default:
			if subject == nil {
				subject = NewSubject()
				obsv := op.notifier(subject.Observable)
				obsv.Subscribe(ctx, func(t Notification) {
					switch {
					case t.HasValue:
						sourceCancel()

						sourceCtx, sourceCancel = context.WithCancel(ctx)
						source.Subscribe(sourceCtx, observer)

					default:
						t.Observe(ob)
						cancel()
					}
				})
			}
			subject.Next(nil)
		}
	}

	sourceCtx, sourceCancel = context.WithCancel(ctx)
	source.Subscribe(sourceCtx, observer)

	return ctx, cancel
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// the exception of a Complete. If the source Observable calls Complete, this
// method will emit to the Observable returned from notifier. If that
// Observable calls Complete or Error, then this method will call Complete or
// Error on the child subscription. Otherwise this method will resubscribe to
// the source Observable.
func (o Observable) RepeatWhen(notifier func(Observable) Observable) Observable {
	op := repeatWhenOperator{notifier}
	return o.Lift(op.Call).Mutex()
}
