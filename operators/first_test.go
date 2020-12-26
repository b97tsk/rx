package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFirst(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.First(),
		),
		rx.ErrEmpty,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.First(),
		),
		ErrTest,
	).Case(
		rx.Just(1).Pipe(
			operators.First(),
		),
		1, Completed,
	).Case(
		rx.Just(1, 2).Pipe(
			operators.First(),
		),
		1, Completed,
	).Case(
		rx.Concat(
			rx.Just(1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.First(),
		),
		1, Completed,
	).Case(
		rx.Concat(
			rx.Just(1, 2),
			rx.Throw(ErrTest),
		).Pipe(
			operators.First(),
		),
		1, Completed,
	).TestAll()
}

func TestFirstOrDefault(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.FirstOrDefault(404),
		),
		404, Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.FirstOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Just(1).Pipe(
			operators.FirstOrDefault(404),
		),
		1, Completed,
	).Case(
		rx.Just(1, 2).Pipe(
			operators.FirstOrDefault(404),
		),
		1, Completed,
	).Case(
		rx.Concat(
			rx.Just(1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.FirstOrDefault(404),
		),
		1, Completed,
	).Case(
		rx.Concat(
			rx.Just(1, 2),
			rx.Throw(ErrTest),
		).Pipe(
			operators.FirstOrDefault(404),
		),
		1, Completed,
	).TestAll()
}
