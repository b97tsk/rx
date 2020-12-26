package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestEvery(t *testing.T) {
	lessThanFive := func(val interface{}, idx int) bool {
		return val.(int) < 5
	}

	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			operators.Every(lessThanFive),
		),
		false, Completed,
	).Case(
		rx.Range(1, 5).Pipe(
			operators.Every(lessThanFive),
		),
		true, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Every(lessThanFive),
		),
		true, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Every(lessThanFive),
		),
		false, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Every(lessThanFive),
		),
		ErrTest,
	).TestAll()
}
