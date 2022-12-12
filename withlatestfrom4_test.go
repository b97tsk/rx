package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom4(t *testing.T) {
	t.Parallel()

	toString := func(v0, v1, v2, v3, v4 string) string {
		return fmt.Sprintf("[%v %v %v %v %v]", v0, v1, v2, v3, v4)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Pipe(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 5)),
			rx.WithLatestFrom4(
				rx.Pipe(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 5)),
				rx.Pipe(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 5)),
				rx.Pipe(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 5)),
				rx.Pipe(rx.Just("E1", "E2", "E3"), AddLatencyToValues[string](5, 5)),
				toString,
			),
		),
		"[A2 B1 C1 D1 E1]",
		"[A3 B2 C2 D2 E2]",
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom4(
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				toString,
			),
		),
		ErrTest,
	)
}
