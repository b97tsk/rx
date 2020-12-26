package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestLast(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Last(),
		),
		rx.ErrEmpty,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Last(),
		),
		ErrTest,
	).Case(
		rx.Just(1).Pipe(
			operators.Last(),
		),
		1, Completed,
	).Case(
		rx.Just(1, 2).Pipe(
			operators.Last(),
		),
		2, Completed,
	).Case(
		rx.Concat(
			rx.Just(1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Last(),
		),
		ErrTest,
	).Case(
		rx.Concat(
			rx.Just(1, 2),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Last(),
		),
		ErrTest,
	).TestAll()
}

func TestLastOrDefault(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.LastOrDefault(404),
		),
		404, Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.LastOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Just(1).Pipe(
			operators.LastOrDefault(404),
		),
		1, Completed,
	).Case(
		rx.Just(1, 2).Pipe(
			operators.LastOrDefault(404),
		),
		2, Completed,
	).Case(
		rx.Concat(
			rx.Just(1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.LastOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Concat(
			rx.Just(1, 2),
			rx.Throw(ErrTest),
		).Pipe(
			operators.LastOrDefault(404),
		),
		ErrTest,
	).TestAll()
}
