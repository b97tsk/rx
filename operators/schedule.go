package operators

import (
	"context"
	"time"
)

func schedule(ctx context.Context, period time.Duration, work func()) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	go func() {
		if period > 0 {
			ticker := time.NewTicker(period)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					for ; !t.After(time.Now()); t = t.Add(period) {
						select {
						case <-done:
							return
						default:
							work()
						}
					}
				}
			}
		} else {
			timer := time.NewTimer(0)
			defer timer.Stop()
			for {
				select {
				case <-done:
					return
				case <-timer.C:
					work()
					timer.Reset(0)
				}
			}
		}
	}()

	return ctx, cancel
}

func scheduleOnce(ctx context.Context, delay time.Duration, work func()) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-done:
		case <-timer.C:
			work()
			cancel()
		}
	}()

	return ctx, cancel
}
