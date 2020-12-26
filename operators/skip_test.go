package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSkip(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 7).Pipe(
			operators.Skip(0),
		),
		1, 2, 3, 4, 5, 6, Completed,
	).Case(
		rx.Range(1, 7).Pipe(
			operators.Skip(3),
		),
		4, 5, 6, Completed,
	).Case(
		rx.Range(1, 3).Pipe(
			operators.Skip(3),
		),
		Completed,
	).Case(
		rx.Range(1, 1).Pipe(
			operators.Skip(3),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 7),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Skip(0),
		),
		1, 2, 3, 4, 5, 6, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 7),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Skip(3),
		),
		4, 5, 6, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 3),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Skip(3),
		),
		ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Skip(3),
		),
		ErrTest,
	).TestAll()
}
