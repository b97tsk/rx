package rx

import (
	"context"
	"time"
)

// A Scheduler is for scheduling some work to be executed asynchronously,
// after a specified time or periodically.
type Scheduler interface {
	Now() time.Time
	Schedule(context.Context, time.Duration, func()) (context.Context, context.CancelFunc)
	ScheduleOnce(context.Context, time.Duration, func()) (context.Context, context.CancelFunc)
}

// DefaultScheduler is used when the optional Scheduler is not specified.
var DefaultScheduler Scheduler = timeScheduler{}
