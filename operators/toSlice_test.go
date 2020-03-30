package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestToSlice(t *testing.T) {
	observables := [...]rx.Observable{
		rx.Just("A", "B", "C"),
		rx.Just("A"),
		rx.Empty(),
		rx.Throw(ErrTest),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(
			operators.ToSlice(),
			operators.Single(),
			ToString(),
		)
	}
	SubscribeN(
		t,
		observables[:],
		[][]interface{}{
			{"[A B C]", rx.Complete},
			{"[A]", rx.Complete},
			{"[]", rx.Complete},
			{ErrTest},
		},
	)
}
