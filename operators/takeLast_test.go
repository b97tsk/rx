package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTakeLast(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			operators.TakeLast(0),
		),
		Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.TakeLast(3),
		),
		6, 7, 8, Completed,
	).Case(
		rx.Range(1, 3).Pipe(
			operators.TakeLast(3),
		),
		1, 2, Completed,
	).Case(
		rx.Range(1, 1).Pipe(
			operators.TakeLast(3),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.TakeLast(0),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.TakeLast(3),
		),
		ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 3),
			rx.Throw(ErrTest),
		).Pipe(
			operators.TakeLast(3),
		),
		ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.TakeLast(3),
		),
		ErrTest,
	).TestAll()
}
