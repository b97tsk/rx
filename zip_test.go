package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/testing"
)

func TestZip(t *testing.T) {
	delay := DelaySubscription(1)
	observables := [...]rx.Observable{
		rx.Zip(rx.Just("A", "B"), rx.Range(1, 4)),
		rx.Zip(rx.Just("A", "B", "C"), rx.Range(1, 4)),
		rx.Zip(rx.Just("A", "B", "C", "D"), rx.Range(1, 4)),
		rx.Zip(rx.Just("A", "B"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
		rx.Zip(rx.Just("A", "B", "C"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
		rx.Zip(rx.Just("A", "B", "C", "D"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(ToString())
	}
	SubscribeN(
		t,
		observables[:],
		[][]interface{}{
			{"[A 1]", "[B 2]", rx.Completed},
			{"[A 1]", "[B 2]", "[C 3]", rx.Completed},
			{"[A 1]", "[B 2]", "[C 3]", rx.Completed},
			{"[A 1]", "[B 2]", rx.Completed},
			{"[A 1]", "[B 2]", "[C 3]", rx.Completed},
			{"[A 1]", "[B 2]", "[C 3]", ErrTest},
		},
	)
	Subscribe(t, rx.Zip(), rx.Completed)
}
