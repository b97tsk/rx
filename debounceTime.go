package rx

import (
	"context"
	"time"
)

type debounceTimeOperator struct {
	source    Operator
	duration  time.Duration
	scheduler Scheduler
}

func (op debounceTimeOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case schedulerOption:
			op.scheduler = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op debounceTimeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCancel := noopFunc
	try := cancellableLocker{}
	latestValue := interface{}(nil)

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = op.scheduler.ScheduleOnce(ctx, op.duration, func() {
			if try.Lock() {
				ob.Next(latestValue)
				try.Unlock()
			}
		})
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				try.Unlock()
				doSchedule()
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.CancelAndUnlock()
				ob.Complete()
				cancel()
			}
		}
	}))

	return ctx, cancel
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
//
// It's like Delay, but passes only the most recent value from each burst of
// emissions.
func (o Observable) DebounceTime(duration time.Duration) Observable {
	op := debounceTimeOperator{
		source:    o.Op,
		duration:  duration,
		scheduler: DefaultScheduler,
	}
	return Observable{op}
}
