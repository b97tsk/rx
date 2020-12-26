package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestCount(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.Count(),
		),
		0, Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.Count(),
		),
		8, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Count(),
		),
		ErrTest,
	).TestAll()
}
