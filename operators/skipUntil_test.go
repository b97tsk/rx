package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSkipUntil(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Just(42)),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Empty()),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Throw(ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Never()),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Just(42).Pipe(DelaySubscription(3))),
		),
		"C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Empty().Pipe(DelaySubscription(3))),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(0, 2),
			operators.SkipUntil(rx.Throw(ErrTest).Pipe(DelaySubscription(3))),
		),
		ErrTest,
	).TestAll()
}
