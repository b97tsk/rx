package rx

import (
	"context"
	"time"
)

type debounceTimeOperator struct {
	Duration time.Duration
}

func (op debounceTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		scheduleCancel = nothingToDo

		latestValue    interface{}
		hasLatestValue bool

		try cancellableLocker
	)

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = scheduleOnce(ctx, op.Duration, func() {
			if try.Lock() {
				defer try.Unlock()
				if hasLatestValue {
					sink.Next(latestValue)
					hasLatestValue = false
				}
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				hasLatestValue = true
				try.Unlock()
				doSchedule()

			case t.HasError:
				try.CancelAndUnlock()
				sink(t)

			default:
				try.CancelAndUnlock()
				if hasLatestValue {
					sink.Next(latestValue)
				}
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
//
// It's like Delay, but passes only the most recent value from each burst of
// emissions.
func (Operators) DebounceTime(duration time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := debounceTimeOperator{duration}
		return source.Lift(op.Call)
	}
}
