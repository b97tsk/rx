package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTimeout(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 1),
			rx.Timeout[string](Step(2)),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.Timeout[string](Step(2)),
		),
		"A", rx.ErrTimeout,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](2, 2),
			rx.Timeout[string](Step(1)).WithFirst(Step(3)),
		),
		"A", rx.ErrTimeout,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.Timeout[string](Step(2)).WithObservable(rx.Throw[string](ErrTest)),
		),
		"A", ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.Timeout[string](Step(2)).WithObservable(rx.Oops[string](ErrTest)),
		),
		"A", rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe3(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.Timeout[string](Step(2)).WithObservable(rx.Oops[string](ErrTest)),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
