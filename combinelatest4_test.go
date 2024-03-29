package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest4(t *testing.T) {
	t.Parallel()

	mapping := func(v1, v2, v3, v4 string) string {
		return fmt.Sprintf("[%v %v %v %v]", v1, v2, v3, v4)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest4(
			rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 4)),
			rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 4)),
			rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 4)),
			rx.Pipe1(rx.Just("D1", "D2"), AddLatencyToValues[string](4, 4)),
			mapping,
		),
		"[A1 B1 C1 D1]",
		"[A2 B1 C1 D1]",
		"[A2 B2 C1 D1]",
		"[A2 B2 C2 D1]",
		"[A2 B2 C2 D2]",
		ErrComplete,
	).Case(
		rx.CombineLatest4(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			mapping,
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.CombineLatest4(
				rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 4)),
				rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 4)),
				rx.Pipe1(rx.Just("D1", "D2"), AddLatencyToValues[string](4, 4)),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
