package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestElementAt(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			operators.ElementAt(4),
		),
		5, Completed,
	).Case(
		rx.Range(1, 5).Pipe(
			operators.ElementAt(4),
		),
		rx.ErrOutOfRange,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.ElementAt(4),
		),
		5, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.ElementAt(4),
		),
		ErrTest,
	).TestAll()
}

func TestElementAtOrDefault(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			operators.ElementAtOrDefault(4, 404),
		),
		5, Completed,
	).Case(
		rx.Range(1, 5).Pipe(
			operators.ElementAtOrDefault(4, 404),
		),
		404, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.ElementAtOrDefault(4, 404),
		),
		5, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.ElementAtOrDefault(4, 404),
		),
		ErrTest,
	).TestAll()
}
