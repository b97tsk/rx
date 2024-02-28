package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest5(t *testing.T) {
	t.Parallel()

	proj := func(v1, v2, v3, v4, v5 string) string {
		return fmt.Sprintf("[%v %v %v %v %v]", v1, v2, v3, v4, v5)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest5(
			rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 5)),
			rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 5)),
			rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 5)),
			rx.Pipe1(rx.Just("D1", "D2"), AddLatencyToValues[string](4, 5)),
			rx.Pipe1(rx.Just("E1", "E2"), AddLatencyToValues[string](5, 5)),
			proj,
		),
		"[A1 B1 C1 D1 E1]",
		"[A2 B1 C1 D1 E1]",
		"[A2 B2 C1 D1 E1]",
		"[A2 B2 C2 D1 E1]",
		"[A2 B2 C2 D2 E1]",
		"[A2 B2 C2 D2 E2]",
		ErrComplete,
	).Case(
		rx.CombineLatest5(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			proj,
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.CombineLatest5(
				rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 5)),
				rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 5)),
				rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("D1", "D2"), AddLatencyToValues[string](4, 5)),
				rx.Pipe1(rx.Just("E1", "E2"), AddLatencyToValues[string](5, 5)),
				proj,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
