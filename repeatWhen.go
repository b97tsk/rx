package rx

import (
	"context"
)

type repeatWhenOperator struct {
	source   Operator
	notifier func(Observable) Observable
}

func (op repeatWhenOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sourceCtx, sourceCancel := canceledCtx, noopFunc

	var subject *Subject
	var observer Observer

	observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			if subject == nil {
				subject = NewSubject()
				obsv := op.notifier(subject.Observable)
				obsv.Subscribe(ctx, ObserverFunc(func(t Notification) {
					switch {
					case t.HasValue:
						sourceCancel()

						sourceCtx, sourceCancel = context.WithCancel(ctx)
						op.source.Call(sourceCtx, observer)

					case t.HasError:
						ob.Error(t.Value.(error))
						cancel()

					default:
						ob.Complete()
						cancel()
					}
				}))
			}
			subject.Next(nil)
		}
	})

	sourceCtx, sourceCancel = context.WithCancel(ctx)
	op.source.Call(sourceCtx, observer)

	return ctx, cancel
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// the exception of a Complete. If the source Observable calls Complete, this
// method will emit to the Observable returned from notifier. If that
// Observable calls Complete or Error, then this method will call Complete or
// Error on the child subscription. Otherwise this method will resubscribe to
// the source Observable.
func (o Observable) RepeatWhen(notifier func(Observable) Observable) Observable {
	op := repeatWhenOperator{
		source:   o.Op,
		notifier: notifier,
	}
	return Observable{op}.Mutex()
}
