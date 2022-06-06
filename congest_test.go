package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCongest(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe4(
			rx.Range(1, 10),
			AddLatencyToValues[int](1, 1),
			rx.Congest[int](0),
			rx.Congest[int](3),
			AddLatencyToValues[int](3, 4),
		),
		1, 2, 3, 4, 5, 6, 7, 8, 9, ErrCompleted,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Pipe2(
			rx.Timer(Step(2)),
			rx.MapTo[time.Time](42),
			rx.Congest[int](3),
		),
		context.DeadlineExceeded,
	)
}
