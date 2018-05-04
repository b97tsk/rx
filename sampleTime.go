package rx

import (
	"context"
	"time"
)

type sampleTimeOperator struct {
	source    Operator
	interval  time.Duration
	scheduler Scheduler
}

func (op sampleTimeOperator) ApplyOptions(options []Option) Operator {
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

func (op sampleTimeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	try := cancellableLocker{}
	latestValue := interface{}(nil)
	hasLatestValue := false

	op.scheduler.Schedule(ctx, op.interval, func() {
		if try.Lock() {
			if hasLatestValue {
				ob.Next(latestValue)
				hasLatestValue = false
			}
			try.Unlock()
		}
	})

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				hasLatestValue = true
				try.Unlock()
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

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func (o Observable) SampleTime(interval time.Duration) Observable {
	op := sampleTimeOperator{
		source:    o.Op,
		interval:  interval,
		scheduler: DefaultScheduler,
	}
	return Observable{op}
}
