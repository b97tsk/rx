package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom7(t *testing.T) {
	t.Parallel()

	proj := func(v0, v1, v2, v3, v4, v5, v6, v7 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v %v]", v0, v1, v2, v3, v4, v5, v6, v7)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 8)),
			rx.WithLatestFrom7(
				rx.Pipe1(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 8)),
				rx.Pipe1(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 8)),
				rx.Pipe1(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 8)),
				rx.Pipe1(rx.Just("E1", "E2", "E3"), AddLatencyToValues[string](5, 8)),
				rx.Pipe1(rx.Just("F1", "F2", "F3"), AddLatencyToValues[string](6, 8)),
				rx.Pipe1(rx.Just("G1", "G2", "G3"), AddLatencyToValues[string](7, 8)),
				rx.Pipe1(rx.Just("H1", "H2", "H3"), AddLatencyToValues[string](8, 8)),
				proj,
			),
		),
		"[A2 B1 C1 D1 E1 F1 G1 H1]",
		"[A3 B2 C2 D2 E2 F2 G2 H2]",
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom7(
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				proj,
			),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Pipe1(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 8)),
			rx.WithLatestFrom7(
				rx.Pipe1(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 8)),
				rx.Pipe1(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 8)),
				rx.Pipe1(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 8)),
				rx.Pipe1(rx.Just("E1", "E2", "E3"), AddLatencyToValues[string](5, 8)),
				rx.Pipe1(rx.Just("F1", "F2", "F3"), AddLatencyToValues[string](6, 8)),
				rx.Pipe1(rx.Just("G1", "G2", "G3"), AddLatencyToValues[string](7, 8)),
				rx.Pipe1(rx.Just("H1", "H2", "H3"), AddLatencyToValues[string](8, 8)),
				proj,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
