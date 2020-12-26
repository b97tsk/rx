package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMap(t *testing.T) {
	double := func(val interface{}, idx int) interface{} {
		return val.(int) * 2
	}

	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Map(double),
		),
		Completed,
	).Case(
		rx.Range(1, 5).Pipe(
			operators.Map(double),
		),
		2, 4, 6, 8, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Map(double),
		),
		2, 4, 6, 8, ErrTest,
	).TestAll()
}

func TestMapTo(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.MapTo(42),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.MapTo(42),
		),
		42, 42, 42, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.MapTo(42),
		),
		42, 42, 42, ErrTest,
	).TestAll()
}
