package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestRepeat(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Repeat[string](0),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Repeat[string](1),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Repeat[string](2),
		),
		"A", "B", "C", "A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Repeat[string](0),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Repeat[string](1),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Repeat[string](2),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.RepeatForever[string](),
		),
		"A", "B", "C", ErrTest,
	)
}
