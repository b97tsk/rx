package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestEndWith(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.EndWith("D", "E"),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.EndWith("D", "E"),
		),
		"D", "E", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.EndWith("D", "E"),
		),
		ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.EndWith(),
		),
		ErrTest,
	).TestAll()
}
