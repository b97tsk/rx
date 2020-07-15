package rx

import (
	"context"
	"time"
)

// Ticker creates an Observable that emits time.Time values every specified
// interval of time.
func Ticker(d time.Duration) Observable {
	return func(ctx context.Context, sink Observer) {
		go func() {
			ticker := time.NewTicker(d)
			defer ticker.Stop()
			done := ctx.Done()
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					sink.Next(t)
				}
			}
		}()
	}
}

// Timer creates an Observable that emits only a time.Time value after
// a particular time span has passed.
func Timer(d time.Duration) Observable {
	return func(ctx context.Context, sink Observer) {
		go func() {
			timer := time.NewTimer(d)
			defer timer.Stop()
			select {
			case <-ctx.Done():
			case t := <-timer.C:
				sink.Next(t)
				sink.Complete()
			}
		}()
	}
}
