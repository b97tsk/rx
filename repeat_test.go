package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestRepeat(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Repeat[string](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Repeat[string](1),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Repeat[string](2),
		),
		"A", "B", "C", "A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Repeat[string](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Repeat[string](1),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Repeat[string](2),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.RepeatForever[string](),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.Repeat[string](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.Repeat[string](1),
		),
		"A", "B", "C", rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.Repeat[string](2),
		),
		"A", "B", "C", rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.RepeatForever[string](),
		),
		"A", "B", "C", rx.ErrOops, ErrTest,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.DoOnComplete[string](func() { time.Sleep(Step(2)) }),
			rx.Repeat[string](2),
		),
		ErrTest,
	)
}
