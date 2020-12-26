package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRetry(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 4).Pipe(
			operators.Retry(0),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Range(1, 4).Pipe(
			operators.Retry(1),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Range(1, 4).Pipe(
			operators.Retry(2),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Retry(0),
		),
		1, 2, 3, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Retry(1),
		),
		1, 2, 3, 1, 2, 3, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Retry(2),
		),
		1, 2, 3, 1, 2, 3, 1, 2, 3, ErrTest,
	).TestAll()
}
