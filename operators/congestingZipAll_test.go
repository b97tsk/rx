package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_CongestingZipAll(t *testing.T) {
	delay := delaySubscription(1)
	observables := [...]Observable{
		Just(Just("A", "B"), Range(1, 4)),
		Just(Just("A", "B", "C"), Range(1, 4)),
		Just(Just("A", "B", "C", "D"), Range(1, 4)),
		Just(Just("A", "B"), Concat(Range(1, 4), Throw(errTest)).Pipe(delay)),
		Just(Just("A", "B", "C"), Concat(Range(1, 4), Throw(errTest)).Pipe(delay)),
		Just(Just("A", "B", "C", "D"), Concat(Range(1, 4), Throw(errTest)).Pipe(delay)),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(operators.CongestingZipAll(), toString)
	}
	subscribeN(
		t,
		observables[:],
		[][]interface{}{
			{"[A 1]", "[B 2]", Complete},
			{"[A 1]", "[B 2]", "[C 3]", Complete},
			{"[A 1]", "[B 2]", "[C 3]", Complete},
			{"[A 1]", "[B 2]", Complete},
			{"[A 1]", "[B 2]", "[C 3]", Complete},
			{"[A 1]", "[B 2]", "[C 3]", errTest},
		},
	)
}
