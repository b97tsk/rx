package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom3(t *testing.T) {
	t.Parallel()

	mapping := func(v0, v1, v2, v3 string) string {
		return fmt.Sprintf("[%v %v %v %v]", v0, v1, v2, v3)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 4)),
			rx.WithLatestFrom3(
				rx.Pipe1(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 4)),
				rx.Pipe1(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 4)),
				mapping,
			),
		),
		"[A2 B1 C1 D1]",
		"[A3 B2 C2 D2]",
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom3(
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
				mapping,
			),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Pipe1(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 4)),
			rx.WithLatestFrom3(
				rx.Pipe1(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("C1", "C2", "C3"), AddLatencyToValues[string](3, 4)),
				rx.Pipe1(rx.Just("D1", "D2", "D3"), AddLatencyToValues[string](4, 4)),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
