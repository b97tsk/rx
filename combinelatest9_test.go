package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest9(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6, v7, v8, v9 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7, v8, v9)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest9(
			rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 9)),
			rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 9)),
			rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 9)),
			rx.Pipe1(rx.Just("D1", "D2"), AddLatencyToValues[string](4, 9)),
			rx.Pipe1(rx.Just("E1", "E2"), AddLatencyToValues[string](5, 9)),
			rx.Pipe1(rx.Just("F1", "F2"), AddLatencyToValues[string](6, 9)),
			rx.Pipe1(rx.Just("G1", "G2"), AddLatencyToValues[string](7, 9)),
			rx.Pipe1(rx.Just("H1", "H2"), AddLatencyToValues[string](8, 9)),
			rx.Pipe1(rx.Just("I1", "I2"), AddLatencyToValues[string](9, 9)),
			toString,
		),
		"[A1 B1 C1 D1 E1 F1 G1 H1 I1]",
		"[A2 B1 C1 D1 E1 F1 G1 H1 I1]",
		"[A2 B2 C1 D1 E1 F1 G1 H1 I1]",
		"[A2 B2 C2 D1 E1 F1 G1 H1 I1]",
		"[A2 B2 C2 D2 E1 F1 G1 H1 I1]",
		"[A2 B2 C2 D2 E2 F1 G1 H1 I1]",
		"[A2 B2 C2 D2 E2 F2 G1 H1 I1]",
		"[A2 B2 C2 D2 E2 F2 G2 H1 I1]",
		"[A2 B2 C2 D2 E2 F2 G2 H2 I1]",
		"[A2 B2 C2 D2 E2 F2 G2 H2 I2]",
		ErrComplete,
	).Case(
		rx.CombineLatest9(
			rx.Throw[string](ErrTest),
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
}
