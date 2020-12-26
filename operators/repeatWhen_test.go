package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRepeatWhen(t *testing.T) {
	repeatOnce := operators.Take(0)
	repeatTwice := operators.Take(1)

	NewTestSuite(t).Case(
		rx.Range(1, 4).Pipe(
			operators.RepeatWhen(repeatOnce),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Range(1, 4).Pipe(
			operators.RepeatWhen(repeatTwice),
		),
		1, 2, 3, 1, 2, 3, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RepeatWhen(repeatOnce),
		),
		1, 2, 3, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RepeatWhen(repeatTwice),
		),
		1, 2, 3, ErrTest,
	).TestAll()
}
