package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFlat(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.Flat(rx.Concat[string]),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.Flat(rx.Concat[string]),
		),
		ErrTest,
	)
}
