package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestStartWith(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("D", "E"),
			rx.StartWith("A", "B", "C"),
		),
		"A", "B", "C", "D", "E", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[string](),
			rx.StartWith("A", "B", "C"),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.StartWith("A", "B", "C"),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.StartWith[string](),
		),
		ErrTest,
	)
}
