package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom2(t *testing.T) {
	t.Parallel()

	toString := func(v0, v1, v2 string) string {
		return fmt.Sprintf("[%v %v %v]", v0, v1, v2)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Pipe(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 3)),
			rx.WithLatestFrom2(
				rx.Pipe(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 3)),
				rx.Pipe(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 3)),
				toString,
			),
		),
		"[A2 B1 C1]",
		"[A3 B2 C2]",
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom2(
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				toString,
			),
		),
		ErrTest,
	)
}
