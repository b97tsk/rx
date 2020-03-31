package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/schedule"
)

// SubscribeOn creates an Observable that asynchronously subscribes to the
// source Observable after waits for the duration to elapse.
func SubscribeOn(d time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return rx.Create(
			func(ctx context.Context, sink rx.Observer) {
				schedule.ScheduleOnce(ctx, d, func() {
					source.Subscribe(ctx, sink)
				})
			},
		)
	}
}
