package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestReduce(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}

		return val
	}

	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}

	NewTestSuite(t).Case(
		rx.Range(1, 7).Pipe(
			operators.Reduce(max),
		),
		6, Completed,
	).Case(
		rx.Just(42).Pipe(
			operators.Reduce(max),
		),
		42, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Reduce(max),
		),
		Completed,
	).Case(
		rx.Range(1, 7).Pipe(
			operators.Reduce(sum),
		),
		21, Completed,
	).Case(
		rx.Just(42).Pipe(
			operators.Reduce(sum),
		),
		42, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Reduce(sum),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Reduce(sum),
		),
		ErrTest,
	).TestAll()
}

func TestFold(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}

		return val
	}

	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}

	NewTestSuite(t).Case(
		rx.Range(1, 7).Pipe(
			operators.Fold(-1, max),
		),
		6, Completed,
	).Case(
		rx.Just(42).Pipe(
			operators.Fold(-1, max),
		),
		42, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Fold(-1, max),
		),
		-1, Completed,
	).Case(
		rx.Range(1, 7).Pipe(
			operators.Fold(-1, sum),
		),
		20, Completed,
	).Case(
		rx.Just(42).Pipe(
			operators.Fold(-1, sum),
		),
		41, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Fold(-1, sum),
		),
		-1, Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Fold(-1, sum),
		),
		ErrTest,
	).TestAll()
}
