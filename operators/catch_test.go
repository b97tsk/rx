package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestCatch(t *testing.T) {
	f := func(error) rx.Observable {
		return rx.Just("D", "E")
	}

	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.Catch(f),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Catch(f),
		),
		"A", "B", "C", "D", "E", Completed,
	).TestAll()
}
