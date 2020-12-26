package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDematerialize(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Materialize(),
			operators.Dematerialize(),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Materialize(),
			operators.Dematerialize(),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Materialize(),
			operators.Dematerialize(),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Materialize(),
			operators.Dematerialize(),
		),
		"A", "B", "C", ErrTest,
	).TestAll()

	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Dematerialize(),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Dematerialize(),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Dematerialize(),
		),
		rx.ErrNotNotification,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Dematerialize(),
		),
		rx.ErrNotNotification,
	).TestAll()
}
