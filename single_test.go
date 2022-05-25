package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSingle(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A"),
			rx.Single[string](),
		),
		"A", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A", "B"),
			rx.Single[string](),
		),
		rx.ErrNotSingle,
	).Case(
		rx.Pipe(
			rx.Empty[string](),
			rx.Single[string](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.Single[string](),
		),
		ErrTest,
	)
}