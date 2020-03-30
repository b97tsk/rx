package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestZip(t *testing.T) {
	delay := delaySubscription(1)
	observables := [...]Observable{
		Zip(Just("A", "B"), Range(1, 4)),
		Zip(Just("A", "B", "C"), Range(1, 4)),
		Zip(Just("A", "B", "C", "D"), Range(1, 4)),
		Zip(Just("A", "B"), Concat(Range(1, 4), Throw(errTest)).Pipe(delay)),
		Zip(Just("A", "B", "C"), Concat(Range(1, 4), Throw(errTest)).Pipe(delay)),
		Zip(Just("A", "B", "C", "D"), Concat(Range(1, 4), Throw(errTest)).Pipe(delay)),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(toString)
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
	subscribe(t, Zip(), Complete)
}
