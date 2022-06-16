package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestRetry(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Retry[string](0),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Retry[string](1),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Retry[string](2),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.RetryForever[string](),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Retry[string](0),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Retry[string](1),
		),
		"A", "B", "C", "A", "B", "C", ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Retry[string](2),
		),
		"A", "B", "C", "A", "B", "C", "A", "B", "C", ErrTest,
	)
}
