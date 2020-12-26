package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSome(t *testing.T) {
	greaterThanFour := func(val interface{}, idx int) bool {
		return val.(int) > 4
	}

	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			operators.Some(greaterThanFour),
		),
		true, Completed,
	).Case(
		rx.Range(1, 5).Pipe(
			operators.Some(greaterThanFour),
		),
		false, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Some(greaterThanFour),
		),
		false, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Some(greaterThanFour),
		),
		true, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Some(greaterThanFour),
		),
		ErrTest,
	).TestAll()
}
