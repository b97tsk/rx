package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest8(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6, v7, v8 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7, v8)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest8(
			rx.Pipe(rx.Just("A", "I"), AddLatencyToValues[string](1, 8)),
			rx.Pipe(rx.Just("B", "J"), AddLatencyToValues[string](2, 8)),
			rx.Pipe(rx.Just("C", "K"), AddLatencyToValues[string](3, 8)),
			rx.Pipe(rx.Just("D", "L"), AddLatencyToValues[string](4, 8)),
			rx.Pipe(rx.Just("E", "M"), AddLatencyToValues[string](5, 8)),
			rx.Pipe(rx.Just("F", "N"), AddLatencyToValues[string](6, 8)),
			rx.Pipe(rx.Just("G", "O"), AddLatencyToValues[string](7, 8)),
			rx.Pipe(rx.Just("H", "P"), AddLatencyToValues[string](8, 8)),
			toString,
		),
		"[A B C D E F G H]", "[I B C D E F G H]", "[I J C D E F G H]",
		"[I J K D E F G H]", "[I J K L E F G H]", "[I J K L M F G H]",
		"[I J K L M N G H]", "[I J K L M N O H]", "[I J K L M N O P]",
		ErrCompleted,
	).Case(
		rx.CombineLatest8(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
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
		rx.CombineLatest8(
			rx.Just("A"),
			rx.Just("B"),
			rx.Just("C"),
			rx.Just("D"),
			rx.Just("E"),
			rx.Just("F"),
			rx.Just("G"),
			rx.Timer(Step(2)),
			func(v1, v2, v3, v4, v5, v6, v7 string, _ time.Time) string {
				return v1 + v2 + v3 + v4 + v5 + v6 + v7
			},
		),
		context.DeadlineExceeded,
	)
}
