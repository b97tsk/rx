package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom8(t *testing.T) {
	t.Parallel()

	toString := func(v0, v1, v2, v3, v4, v5, v6, v7, v8 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v %v %v]", v0, v1, v2, v3, v4, v5, v6, v7, v8)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Pipe(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 9)),
			rx.WithLatestFrom8(
				rx.Pipe(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 9)),
				rx.Pipe(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 9)),
				rx.Pipe(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 9)),
				rx.Pipe(rx.Just("E1", "E2", "E3"), AddLatencyToValues[string](5, 9)),
				rx.Pipe(rx.Just("F1", "F2", "F3"), AddLatencyToValues[string](6, 9)),
				rx.Pipe(rx.Just("G1", "G2", "G3"), AddLatencyToValues[string](7, 9)),
				rx.Pipe(rx.Just("H1", "H2", "H3"), AddLatencyToValues[string](8, 9)),
				rx.Pipe(rx.Just("I1", "I2", "I3"), AddLatencyToValues[string](9, 9)),
				toString,
			),
		),
		"[A2 B1 C1 D1 E1 F1 G1 H1 I1]",
		"[A3 B2 C2 D2 E2 F2 G2 H2 I2]",
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom8(
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
		),
		ErrTest,
	)
}
