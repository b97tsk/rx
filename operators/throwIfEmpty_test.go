package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestThrowIfEmpty(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.ThrowIfEmpty(rx.ErrEmpty),
		),
		rx.ErrEmpty,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.ThrowIfEmpty(rx.ErrEmpty),
		),
		ErrTest,
	).Case(
		rx.Just(1).Pipe(
			operators.ThrowIfEmpty(rx.ErrEmpty),
		),
		1, Completed,
	).Case(
		rx.Just(1, 2).Pipe(
			operators.ThrowIfEmpty(rx.ErrEmpty),
		),
		1, 2, Completed,
	).Case(
		rx.Concat(
			rx.Just(1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.ThrowIfEmpty(rx.ErrEmpty),
		),
		1, ErrTest,
	).Case(
		rx.Concat(
			rx.Just(1, 2),
			rx.Throw(ErrTest),
		).Pipe(
			operators.ThrowIfEmpty(rx.ErrEmpty),
		),
		1, 2, ErrTest,
	).TestAll()
}
