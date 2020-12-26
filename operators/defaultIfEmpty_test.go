package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDefaultIfEmpty(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Empty().Pipe(
			operators.DefaultIfEmpty(42),
		),
		42, Completed,
	).Case(
		rx.Range(1, 4).Pipe(
			operators.DefaultIfEmpty(42),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
		).Pipe(
			operators.DefaultIfEmpty(42),
		),
		1, 2, 3, ErrTest,
	).TestAll()
}
