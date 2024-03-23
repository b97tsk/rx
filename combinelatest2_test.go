package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest2(t *testing.T) {
	t.Parallel()

	mapping := func(v1, v2 string) string {
		return fmt.Sprintf("[%v %v]", v1, v2)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest2(
			rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 2)),
			rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 2)),
			mapping,
		),
		"[A1 B1]",
		"[A2 B1]",
		"[A2 B2]",
		ErrComplete,
	).Case(
		rx.CombineLatest2(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			mapping,
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.CombineLatest2(
				rx.Pipe1(rx.Just("A1", "A2"), AddLatencyToValues[string](1, 2)),
				rx.Pipe1(rx.Just("B1", "B2"), AddLatencyToValues[string](2, 2)),
				mapping,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
