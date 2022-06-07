package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest2(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2 string) string {
		return fmt.Sprintf("[%v %v]", v1, v2)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest2(
			rx.Pipe(rx.Just("A", "C"), AddLatencyToValues[string](1, 2)),
			rx.Pipe(rx.Just("B", "D"), AddLatencyToValues[string](2, 2)),
			toString,
		),
		"[A B]", "[C B]", "[C D]", ErrCompleted,
	).Case(
		rx.CombineLatest2(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			toString,
		),
		ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.CombineLatest2(
			rx.Just("A"),
			rx.Timer(Step(2)),
			func(v1 string, _ time.Time) string {
				return v1
			},
		),
		context.DeadlineExceeded,
	)
}
