package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom5(t *testing.T) {
	t.Parallel()

	toString := func(v0, v1, v2, v3, v4, v5 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v]", v0, v1, v2, v3, v4, v5)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 6)),
			rx.WithLatestFrom5(
				rx.Pipe1(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 6)),
				rx.Pipe1(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 6)),
				rx.Pipe1(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 6)),
				rx.Pipe1(rx.Just("E1", "E2", "E3"), AddLatencyToValues[string](5, 6)),
				rx.Pipe1(rx.Just("F1", "F2", "F3"), AddLatencyToValues[string](6, 6)),
				toString,
			),
		),
		"[A2 B1 C1 D1 E1 F1]",
		"[A3 B2 C2 D2 E2 F2]",
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom5(
				rx.Throw[string](ErrTest),
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
