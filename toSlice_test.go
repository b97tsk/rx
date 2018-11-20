package rx_test

import (
	"fmt"
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
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return fmt.Sprint(val)
				},
			),
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
