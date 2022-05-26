package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestEndWith(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.EndWith("D", "E"),
		),
		"A", "B", "C", "D", "E", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[string](),
			rx.EndWith("D", "E"),
		),
		"D", "E", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.EndWith("D", "E"),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.EndWith[string](),
		),
		ErrTest,
	)
}
