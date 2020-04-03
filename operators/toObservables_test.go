package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestToObservables(t *testing.T) {
	observables := [...]rx.Observable{
		rx.Just("A", "B", "C"),
		rx.Just(rx.Just("A"), rx.Just("B"), rx.Just("C")),
		rx.Empty(),
		rx.Throw(ErrTest),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(
			operators.ToObservables(),
			operators.Single(),
			operators.ConcatMap(
				func(val interface{}, idx int) rx.Observable {
					return rx.Concat(val.([]rx.Observable)...)
				},
			),
		)
	}
	SubscribeN(
		t,
		observables[:],
		[][]interface{}{
			{rx.ErrNotObservable},
			{"A", "B", "C", rx.Complete},
			{rx.Complete},
			{ErrTest},
		},
	)
}