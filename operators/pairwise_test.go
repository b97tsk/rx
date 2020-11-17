package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestPairwise(t *testing.T) {
	observables := [...]rx.Observable{
		rx.Empty(),
		rx.Just("A"),
		rx.Just("A", "B"),
		rx.Just("A", "B", "C"),
		rx.Just("A", "B", "C", "D"),
		rx.Concat(rx.Just("A", "B", "C", "D"), rx.Throw(ErrTest)),
	}

	for i, obs := range observables {
		observables[i] = obs.Pipe(operators.Pairwise(), ToString())
	}

	SubscribeN(
		t,
		observables[:],
		[][]interface{}{
			{Completed},
			{Completed},
			{"{A B}", Completed},
			{"{A B}", "{B C}", Completed},
			{"{A B}", "{B C}", "{C D}", Completed},
			{"{A B}", "{B C}", "{C D}", ErrTest},
		},
	)
}
