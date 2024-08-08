package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDelay(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C", "D", "E"),
			rx.Delay[string](Step(1)),
		),
		"A", "B", "C", "D", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 1),
			rx.Delay[string](Step(2)),
		),
		"A", "B", "C", "D", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Concat(
				rx.Just("A", "B", "C", "D", "E"),
				rx.Throw[string](ErrTest),
			),
			AddLatencyToNotifications[string](0, 2),
			rx.Delay[string](Step(1)),
		),
		"A", "B", "C", "D", "E", ErrTest,
	).Case(
		rx.Pipe2(
			rx.Concat(
				rx.Just("A", "B", "C", "D", "E"),
				rx.Oops[string](ErrTest),
			),
			AddLatencyToNotifications[string](0, 2),
			rx.Delay[string](Step(1)),
		),
		"A", "B", "C", "D", rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.Delay[string](Step(1)),
		),
		ErrComplete,
	)

	{
		ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
		defer cancel()

		NewTestSuite[string](t).WithContext(ctx).Case(
			rx.Pipe1(
				rx.Just("A", "B", "C", "D", "E"),
				rx.Delay[string](Step(2)),
			),
			ErrTest,
		)
	}

	{
		ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(2), ErrTest)
		defer cancel()

		NewTestSuite[string](t).WithContext(ctx).Case(
			rx.Pipe2(
				rx.Just("A", "B", "C", "D", "E"),
				rx.Delay[string](Step(1)),
				rx.DoOnNext(func(string) { time.Sleep(Step(2)) }),
			),
			"A", ErrTest,
		)
	}
}
