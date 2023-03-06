package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestWithLatestFrom1(t *testing.T) {
	t.Parallel()

	toString := func(v0, v1 string) string {
		return fmt.Sprintf("[%v %v]", v0, v1)
	}

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A1", "A2", "A3"), AddLatencyToValues[string](1, 2)),
			rx.WithLatestFrom1(
				rx.Pipe1(rx.Just("B1", "B2", "B3"), AddLatencyToValues[string](2, 2)),
				toString,
			),
		),
		"[A2 B1]",
		"[A3 B2]",
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.WithLatestFrom1(
				rx.Throw[string](ErrTest),
				toString,
			),
		),
		ErrTest,
	)
}
