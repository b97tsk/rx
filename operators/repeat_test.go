package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRepeat(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Repeat(0),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Repeat(1),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Repeat(2),
		),
		"A", "B", "C", "A", "B", "C", Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Repeat(0),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Repeat(1),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Repeat(2),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RepeatForever(),
		),
		"A", "B", "C", ErrTest,
	).TestAll()
}
