package rx

import (
	"context"
	"math"
	"sync/atomic"
)

type retryWhenOperator struct {
	Notifier func(Observable) Observable
}

func (op retryWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var sourceLocker *cancellableLocker
	sourceCtx, sourceCancel := Done()

	var (
		activeCount    = uint32(2)
		lastError      error
		subject        Subject
		createSubject  func() Subject
		avoidRecursive avoidRecursiveCalls
	)

	subscribe := func() {
		var try cancellableLocker
		sourceLocker = &try
		sourceCtx, sourceCancel = context.WithCancel(ctx)
		source.Subscribe(sourceCtx, func(t Notification) {
			if try.Lock() {
				switch {
				case t.HasValue:
					sink(t)
					try.Unlock()
				case t.HasError:
					lastError = t.Value.(error)
					activeCount := atomic.AddUint32(&activeCount, math.MaxUint32)
					try.CancelAndUnlock()
					if activeCount == 0 {
						sink(t)
						break
					}
					if subject.Observer == nil {
						subject = createSubject()
					}
					subject.Next(t.Value)
				default:
					sink(t)
					try.Unlock()
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
				if sourceLocker.Lock() {
					sourceLocker.CancelAndUnlock()
				} else {
					atomic.AddUint32(&activeCount, 1)
				}
				avoidRecursive.Do(subscribe)

			case t.HasError:
				sink(t)

			default:
				if atomic.AddUint32(&activeCount, math.MaxUint32) == 0 {
					sink.Error(lastError)
				}
			}
		})
		return subject
	}

	avoidRecursive.Do(subscribe)

	return ctx, cancel
}

// RetryWhen creates an Observable that mirrors the source Observable with
// the exception of ERROR emission. If the source Observable errors, this
// operator will emit the error to the Observable returned from notifier.
// If that Observable emits a value, this operator will resubscribe to the
// source Observable. Otherwise, this operator will emit the last error on
// the child subscription.
func (Operators) RetryWhen(notifier func(Observable) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := retryWhenOperator{notifier}
		return source.Lift(op.Call)
	}
}
