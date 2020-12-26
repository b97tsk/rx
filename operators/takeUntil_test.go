package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTakeUntil(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Just(42)),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Empty()),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Throw(ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Never()),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Just(42).Pipe(DelaySubscription(3))),
		),
		"A", "B", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Empty().Pipe(DelaySubscription(3))),
		),
		"A", "B", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.TakeUntil(rx.Throw(ErrTest).Pipe(DelaySubscription(3))),
		),
		"A", "B", ErrTest,
	).TestAll()
}
