package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestIsEmpty(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B").Pipe(
			operators.IsEmpty(),
		),
		false, Completed,
	).Case(
		rx.Just("A").Pipe(
			operators.IsEmpty(),
		),
		false, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.IsEmpty(),
		),
		true, Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.IsEmpty(),
		),
		ErrTest,
	).TestAll()
}
