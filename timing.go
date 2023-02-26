package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/timerpool"
)

// Ticker creates an Observable that emits time.Time values every specified
// interval of time.
func Ticker(d time.Duration) Observable[time.Time] {
	if d <= 0 {
		panic("d <= 0")
	}

	return func(ctx context.Context, sink Observer[time.Time]) {
		tk := time.NewTicker(d)

		Go(ctx, func() {
			defer tk.Stop()

			done := ctx.Done()

			for {
				select {
				case <-done:
					sink.Error(ctx.Err())
					return
				case t := <-tk.C:
					sink.Next(t)
				}
			}
		})
	}
}

// Timer creates an Observable that emits only a time.Time value after
// a particular time span has passed.
func Timer(d time.Duration) Observable[time.Time] {
	return func(ctx context.Context, sink Observer[time.Time]) {
		tm := timerpool.Get(d)

		Go(ctx, func() {
			select {
			case <-ctx.Done():
				timerpool.Put(tm)
				sink.Error(ctx.Err())
			case t := <-tm.C:
				timerpool.PutExpired(tm)
				sink.Next(t)
				sink.Complete()
			}
		})
	}
}
