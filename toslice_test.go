package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestToSlice(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		"[A B C]", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A"),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		"[A]", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		"[]", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		ErrTest,
	)
}
