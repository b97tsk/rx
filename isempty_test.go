package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	NewTestSuite[bool](t).Case(
		rx.Pipe(
			rx.Just("A", "B"),
			rx.IsEmpty[string](),
		),
		false, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A"),
			rx.IsEmpty[string](),
		),
		false, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[string](),
			rx.IsEmpty[string](),
		),
		true, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[string](ErrTest),
			rx.IsEmpty[string](),
		),
		ErrTest,
	)
}
