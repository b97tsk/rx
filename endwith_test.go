package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestEndWith(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.EndWith[string](),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.EndWith("D", "E"),
		),
		"A", "B", "C", "D", "E", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.EndWith("D", "E"),
		),
		"D", "E", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.EndWith("D", "E"),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.EndWith("D", "E"),
		),
		rx.ErrOops, ErrTest,
	)
}
