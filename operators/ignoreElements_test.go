package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestIgnoreElements(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.IgnoreElements(),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.IgnoreElements(),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.IgnoreElements(),
		),
		ErrTest,
	).TestAll()
}
