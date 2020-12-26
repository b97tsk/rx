package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSingle(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B").Pipe(
			operators.Single(),
		),
		rx.ErrNotSingle,
	).Case(
		rx.Just("A").Pipe(
			operators.Single(),
		),
		"A", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Single(),
		),
		rx.ErrEmpty,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Single(),
		),
		ErrTest,
	).TestAll()
}
