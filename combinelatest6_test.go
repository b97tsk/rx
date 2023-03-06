package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest6(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest6(
			rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 6)),
			rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 6)),
			rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 6)),
			rx.Pipe1(rx.Just("D1", "D2"), AddLatencyToValues[string](4, 6)),
			rx.Pipe1(rx.Just("E1", "E2"), AddLatencyToValues[string](5, 6)),
			rx.Pipe1(rx.Just("F1", "F2"), AddLatencyToValues[string](6, 6)),
			toString,
		),
		"[A1 B1 C1 D1 E1 F1]",
		"[A2 B1 C1 D1 E1 F1]",
		"[A2 B2 C1 D1 E1 F1]",
		"[A2 B2 C2 D1 E1 F1]",
		"[A2 B2 C2 D2 E1 F1]",
		"[A2 B2 C2 D2 E2 F1]",
		"[A2 B2 C2 D2 E2 F2]",
		ErrComplete,
	).Case(
		rx.CombineLatest6(
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
