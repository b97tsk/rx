package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCatch(t *testing.T) {
	t.Parallel()

	f := func(error) rx.Observable[string] {
		return rx.Just("D", "E")
	}

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.Catch(f),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Catch(f),
		),
		"A", "B", "C", "D", "E", ErrCompleted,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe(
			rx.Never[string](),
			rx.Catch(f),
		),
		context.DeadlineExceeded,
	)
}

func TestOnErrorResumeWith(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.OnErrorResumeWith(rx.Just("D", "E")),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.OnErrorResumeWith(rx.Just("D", "E")),
		),
		"A", "B", "C", "D", "E", ErrCompleted,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe(
			rx.Never[string](),
			rx.OnErrorResumeWith(rx.Just("D", "E")),
		),
		context.DeadlineExceeded,
	)
}

func TestOnErrorComplete(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe(
			rx.Just("A", "B", "C"),
			rx.OnErrorComplete[string](),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.OnErrorComplete[string](),
		),
		"A", "B", "C", ErrCompleted,
	)
}
