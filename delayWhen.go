package rx

import (
	"context"
	"math"
	"sync/atomic"
)

type delayWhenOperator struct {
	DurationSelector func(interface{}, int) Observable
}

func (op delayWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		outerIndex  = -1
		activeCount = uint32(1)
	)

	doSchedule := func(val interface{}, index int) {
		scheduleCtx, scheduleCancel := context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			scheduleCancel()
			switch {
			case t.HasValue:
				sink.Next(val)

				if atomic.AddUint32(&activeCount, math.MaxUint32) == 0 {
					sink.Complete()
				}

			case t.HasError:
				sink(t)

			default:
				if atomic.AddUint32(&activeCount, math.MaxUint32) == 0 {
					sink(t)
				}
			}
		}

		obs := op.DurationSelector(val, index)
		obs.Subscribe(scheduleCtx, observer.Notify)
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++
			atomic.AddUint32(&activeCount, 1)
			doSchedule(t.Value, outerIndex)
		case t.HasError:
			sink(t)
		default:
			if atomic.AddUint32(&activeCount, math.MaxUint32) == 0 {
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
