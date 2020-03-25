package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ToObservables(t *testing.T) {
	observables := [...]Observable{
		Just("A", "B", "C"),
		Just(Just("A"), Just("B"), Just("C")),
		Empty(),
		Throw(errTest),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(
			operators.ToObservables(),
			operators.Single(),
			operators.ConcatMap(
				func(val interface{}, idx int) Observable {
					return Concat(val.([]Observable)...)
				},
			),
		)
	}
	subscribe(
		t, observables[:],
		ErrNotObservable,
		"A", "B", "C", Complete,
		Complete,
		errTest,
	)
}
