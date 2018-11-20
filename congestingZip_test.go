package rx_test

import (
	"fmt"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_CongestingZipAll(t *testing.T) {
	delay := delaySubscription(1)
	observables := [...]Observable{
		Just(Just("A", "B"), Range(1, 4)),
		Just(Just("A", "B", "C"), Range(1, 4)),
		Just(Just("A", "B", "C", "D"), Range(1, 4)),
		Just(Just("A", "B"), Concat(Range(1, 4), Throw(xErrTest)).Pipe(delay)),
		Just(Just("A", "B", "C"), Concat(Range(1, 4), Throw(xErrTest)).Pipe(delay)),
		Just(Just("A", "B", "C", "D"), Concat(Range(1, 4), Throw(xErrTest)).Pipe(delay)),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(
			operators.CongestingZipAll(),
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return fmt.Sprintf("%v%v", val.([]interface{})...)
				},
			),
		)
	}
	subscribe(
		t, observables[:],
		"A1", "B2", xComplete,
		"A1", "B2", "C3", xComplete,
		"A1", "B2", "C3", xComplete,
		"A1", "B2", xComplete,
		"A1", "B2", "C3", xComplete,
		"A1", "B2", "C3", xErrTest,
	)
	subscribe(t, []Observable{CongestingZip()}, xComplete)
}
