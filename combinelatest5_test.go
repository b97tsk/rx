package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest5(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5 string) string {
		return fmt.Sprintf("[%v %v %v %v %v]", v1, v2, v3, v4, v5)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest5(
			rx.Pipe(rx.Just("A", "F"), AddLatencyToValues[string](1, 5)),
			rx.Pipe(rx.Just("B", "G"), AddLatencyToValues[string](2, 5)),
			rx.Pipe(rx.Just("C", "H"), AddLatencyToValues[string](3, 5)),
			rx.Pipe(rx.Just("D", "I"), AddLatencyToValues[string](4, 5)),
			rx.Pipe(rx.Just("E", "J"), AddLatencyToValues[string](5, 5)),
			toString,
		),
		"[A B C D E]", "[F B C D E]", "[F G C D E]",
		"[F G H D E]", "[F G H I E]", "[F G H I J]",
		ErrCompleted,
	).Case(
		rx.CombineLatest5(
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
		rx.CombineLatest5(
			rx.Just("A"),
			rx.Just("B"),
			rx.Just("C"),
			rx.Just("D"),
			rx.Timer(Step(2)),
			func(v1, v2, v3, v4 string, _ time.Time) string {
				return v1 + v2 + v3 + v4
			},
		),
		context.DeadlineExceeded,
	)
}