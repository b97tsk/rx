package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type delayWhenOperator struct {
	DurationSelector func(interface{}, int) Observable
}

func (op delayWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		outerIndex  = -1
		activeCount = atomic.Uint32(1)
	)

	doSchedule := func(val interface{}, idx int) {
		scheduleCtx, scheduleCancel := context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			scheduleCancel()
			switch {
			case t.HasValue:
				sink.Next(val)

				if activeCount.Sub(1) == 0 {
					sink.Complete()
				}

			case t.HasError:
				sink(t)

			default:
				if activeCount.Sub(1) == 0 {
					sink(t)
				}
			}
		}

		obs := op.DurationSelector(val, idx)
		obs.Subscribe(scheduleCtx, observer.Notify)
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++
			activeCount.Add(1)
			doSchedule(t.Value, outerIndex)
		case t.HasError:
			sink(t)
		default:
			if activeCount.Sub(1) == 0 {
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// DelayWhen creates an Observable that delays the emission of items from
// the source Observable by a given time span determined by the emissions of
// another Observable.
//
// It's like Delay, but the time span of the delay duration is determined by
// a second Observable.
func (Operators) DelayWhen(durationSelector func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := delayWhenOperator{durationSelector}
		return source.Lift(op.Call)
	}
}
