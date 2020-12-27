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
		rx.Just("A", "B", "C").Pipe(
			operators.RepeatWhen(repeatOnce),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.RepeatWhen(repeatTwice),
		),
		"A", "B", "C", "A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.RepeatWhen(
				func(rx.Observable) rx.Observable {
					return rx.Throw(ErrTest)
				},
			),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RepeatWhen(repeatOnce),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RepeatWhen(repeatTwice),
		),
		"A", "B", "C", ErrTest,
	).TestAll()
}
