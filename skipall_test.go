package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSkipAll(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.SkipAll[string](),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.SkipAll[string](),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.SkipAll[string](),
		),
		ErrTest,
	)
}
