package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFilter(t *testing.T) {
	lessThanFive := func(val interface{}, idx int) bool {
		return val.(int) < 5
	}

	NewTestSuite(t).Case(
		rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(
			operators.Filter(lessThanFive),
		),
		1, 2, 3, 4, 4, 3, 2, 1, Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.Filter(lessThanFive),
		),
		1, 2, 3, 4, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Filter(lessThanFive),
		),
		1, 2, 3, 4, ErrTest,
	).TestAll()
}

func TestFilterMap(t *testing.T) {
	lessThanFive := func(val interface{}, idx int) (interface{}, bool) {
		return val, val.(int) < 5
	}

	NewTestSuite(t).Case(
		rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(
			operators.FilterMap(lessThanFive),
		),
		1, 2, 3, 4, 4, 3, 2, 1, Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.FilterMap(lessThanFive),
		),
		1, 2, 3, 4, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.FilterMap(lessThanFive),
		),
		1, 2, 3, 4, ErrTest,
	).TestAll()
}

func TestExclude(t *testing.T) {
	lessThanFive := func(val interface{}, idx int) bool {
		return val.(int) < 5
	}

	NewTestSuite(t).Case(
		rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(
			operators.Exclude(lessThanFive),
		),
		5, Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.Exclude(lessThanFive),
		),
		5, 6, 7, 8, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Exclude(lessThanFive),
		),
		5, 6, 7, 8, ErrTest,
	).TestAll()
}
