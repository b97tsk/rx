package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMaterialize(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Materialize(),
			operators.Count(),
		),
		1, Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Materialize(),
			operators.Count(),
		),
		1, Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Materialize(),
			operators.Count(),
		),
		4, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Materialize(),
			operators.Count(),
		),
		4, Completed,
	).TestAll()
}
