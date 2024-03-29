package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest3(t *testing.T) {
	t.Parallel()

	mapping := func(v1, v2, v3 string) string {
		return fmt.Sprintf("[%v %v %v]", v1, v2, v3)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest3(
			rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 3)),
			rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 3)),
			rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 3)),
			mapping,
		),
		"[A1 B1 C1]",
		"[A2 B1 C1]",
		"[A2 B2 C1]",
		"[A2 B2 C2]",
		ErrComplete,
	).Case(
		rx.CombineLatest3(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			mapping,
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.CombineLatest3(
				rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 3)),
				rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 3)),
				rx.Pipe1(rx.Just("C1", "C2"), AddLatencyToValues[string](3, 3)),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
