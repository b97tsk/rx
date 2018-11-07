package rx

import (
	"context"
)

type repeatWhenOperator struct {
	Notifier func(Observable) Observable
}

func (op repeatWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sourceCtx, sourceCancel := canceledCtx, nothingToDo

	sink = Mutex(Finally(sink, cancel))

	var (
		subject        *Subject
		observer       Observer
		avoidRecursive avoidRecursiveCalls
	)

	subscribe := func() {
		sourceCtx, sourceCancel = context.WithCancel(ctx)
		source.Subscribe(sourceCtx, observer)
	}

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			sink(t)
		default:
			if subject == nil {
				subject = NewSubject()
				obsv := op.Notifier(subject.Observable)
				obsv.Subscribe(ctx, func(t Notification) {
					switch {
					case t.HasValue:
						sourceCancel()
						avoidRecursive.Do(subscribe)

					default:
						sink(t)
					}
				})
			}
			subject.Next(nil)
		}
	}

	avoidRecursive.Do(subscribe)

	return ctx, cancel
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// the exception of a Complete. If the source Observable calls Complete, this
// method will emit to the Observable returned from notifier. If that
// Observable calls Complete or Error, then this method will call Complete or
// Error on the child subscription. Otherwise this method will resubscribe to
// the source Observable.
func (Operators) RepeatWhen(notifier func(Observable) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := repeatWhenOperator{notifier}
		return source.Lift(op.Call)
	}
}
