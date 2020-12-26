package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestToSlice(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.ToSlice(),
			operators.Single(),
			ToString(),
		),
		"[A B C]", Completed,
	).Case(
		rx.Just("A").Pipe(
			operators.ToSlice(),
			operators.Single(),
			ToString(),
		),
		"[A]", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.ToSlice(),
			operators.Single(),
			ToString(),
		),
		"[]", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.ToSlice(),
			operators.Single(),
			ToString(),
		),
		ErrTest,
	).TestAll()
}
