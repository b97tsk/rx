package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestStartWith(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("D", "E").Pipe(
			operators.StartWith("A", "B", "C"),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.StartWith("A", "B", "C"),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.StartWith("A", "B", "C"),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.StartWith(),
		),
		ErrTest,
	).TestAll()
}
