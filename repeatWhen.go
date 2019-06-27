package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type repeatWhenOperator struct {
	Notifier func(Observable) Observable
}

func (op repeatWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	sourceCtx, sourceCancel := Done()

	var (
		activeCount    = atomic.Uint32(2)
		subject        Subject
		createSubject  func() Subject
		avoidRecursive avoidRecursiveCalls
	)

	type X struct{}
	var cxCurrent chan X

	subscribe := func() {
		cx := make(chan X, 1)
		cx <- X{}
		cxCurrent = cx
		sourceCtx, sourceCancel = context.WithCancel(ctx)
		source.Subscribe(sourceCtx, func(t Notification) {
			if x, ok := <-cx; ok {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
					cx <- x
				default:
					activeCount := activeCount.Sub(1)
					close(cx)
					if activeCount == 0 {
						sink(t)
						break
					}
					if subject.Observer == nil {
						subject = createSubject()
					}
					subject.Next(nil)
				}
			}
		})
	}

	createSubject = func() Subject {
		subject := NewSubject()
		obs := op.Notifier(subject.Observable)
		obs.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				sourceCancel()
				if _, ok := <-cxCurrent; ok {
					close(cxCurrent)
				} else {
					activeCount.Add(1)
				}
				avoidRecursive.Do(subscribe)

			case t.HasError:
				sink(t)

			default:
				if activeCount.Sub(1) == 0 {
					sink(t)
				}
			}
		})
		return subject
	}

	avoidRecursive.Do(subscribe)

	return ctx, cancel
}

// RepeatWhen creates an Observable that mirrors the source Observable with
// the exception of COMPLETE emission. If the source Observable completes,
// this operator will emit nil to the Observable returned from notifier. If
// that Observable emits a value, this operator will resubscribe to the source
// Observable. Otherwise, this operator will emit a COMPLETE on the child
// subscription.
func (Operators) RepeatWhen(notifier func(Observable) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := repeatWhenOperator{notifier}
		return source.Lift(op.Call)
	}
}
