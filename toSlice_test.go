package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ToSlice(t *testing.T) {
	observables := [...]Observable{
		Just("A", "B", "C"),
		Just("A"),
		Empty(),
		Throw(xErrTest),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(
			operators.ToSlice(),
			operators.Single(),
			toString,
		)
	}
	subscribe(
		t, observables[:],
		"[A B C]", xComplete,
		"[A]", xComplete,
		"[]", xComplete,
		xErrTest,
	)
}