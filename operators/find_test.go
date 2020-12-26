package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFind(t *testing.T) {
	equalsFive := func(val interface{}, idx int) bool {
		return val.(int) == 5
	}

	NewTestSuite(t).Case(
		rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(
			operators.Find(equalsFive),
		),
		5, Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.Find(equalsFive),
		),
		5, Completed,
	).Case(
		rx.Range(1, 5).Pipe(
			operators.Find(equalsFive),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Find(equalsFive),
		),
		5, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Find(equalsFive),
		),
		ErrTest,
	).TestAll()
}
