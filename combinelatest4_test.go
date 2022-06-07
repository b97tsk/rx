package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest4(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4 string) string {
		return fmt.Sprintf("[%v %v %v %v]", v1, v2, v3, v4)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest4(
			rx.Pipe(rx.Just("A", "E"), AddLatencyToValues[string](1, 4)),
			rx.Pipe(rx.Just("B", "F"), AddLatencyToValues[string](2, 4)),
			rx.Pipe(rx.Just("C", "G"), AddLatencyToValues[string](3, 4)),
			rx.Pipe(rx.Just("D", "H"), AddLatencyToValues[string](4, 4)),
			toString,
		),
		"[A B C D]", "[E B C D]", "[E F C D]", "[E F G D]", "[E F G H]", ErrCompleted,
	).Case(
		rx.CombineLatest4(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			toString,
		),
		ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.CombineLatest4(
			rx.Just("A"),
			rx.Just("B"),
			rx.Just("C"),
			rx.Timer(Step(2)),
			func(v1, v2, v3 string, _ time.Time) string {
				return v1 + v2 + v3
			},
		),
		context.DeadlineExceeded,
	)
}
