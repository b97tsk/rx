package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestPairwise(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Pairwise(),
			ToString(),
		),
		Completed,
	).Case(
		rx.Just("A").Pipe(
			operators.Pairwise(),
			ToString(),
		),
		Completed,
	).Case(
		rx.Just("A", "B").Pipe(
			operators.Pairwise(),
			ToString(),
		),
		"{A B}", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Pairwise(),
			ToString(),
		),
		"{A B}", "{B C}", Completed,
	).Case(
		rx.Just("A", "B", "C", "D").Pipe(
			operators.Pairwise(),
			ToString(),
		),
		"{A B}", "{B C}", "{C D}", Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C", "D"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Pairwise(),
			ToString(),
		),
		"{A B}", "{B C}", "{C D}", ErrTest,
	).TestAll()
}
