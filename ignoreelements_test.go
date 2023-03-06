package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestIgnoreElements(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.IgnoreElements[string, string](),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.IgnoreElements[string, string](),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.IgnoreElements[string, string](),
		),
		ErrTest,
	)
}
