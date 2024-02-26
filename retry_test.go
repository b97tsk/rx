package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestRetry(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Retry[string](0),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Retry[string](1),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Retry[string](2),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.RetryForever[string](),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Retry[string](0),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Retry[string](1),
		),
		"A", "B", "C", "A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Retry[string](2),
		),
		"A", "B", "C", "A", "B", "C", "A", "B", "C", ErrTest,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe1(
			rx.Never[string](),
			rx.Retry[string](1),
		),
		context.DeadlineExceeded,
	)
}
