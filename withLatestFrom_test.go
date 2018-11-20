package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_WithLatestFrom(t *testing.T) {
	addLatency1 := addLatencyToValue(1, 2)
	addLatency2 := addLatencyToNotification(0, 2)

	observables := [...]Observable{
		Just("A", "B").Pipe(addLatency1),
		Just("A", "B", "C").Pipe(addLatency1),
		Just("A", "B", "C", "D").Pipe(addLatency1),
	}

	{
		observables := observables
		for i, obs := range observables {
			observables[i] = obs.Pipe(
				operators.WithLatestFrom(Range(1, 4).Pipe(addLatency2)),
				toString,
			)
		}
		subscribe(
			t, observables[:],
			"[A 1]", "[B 2]", xComplete,
			"[A 1]", "[B 2]", "[C 3]", xComplete,
			"[A 1]", "[B 2]", "[C 3]", "[D 3]", xComplete,
		)
	}

	{
		observables := observables
		for i, obs := range observables {
			observables[i] = obs.Pipe(
				operators.WithLatestFrom(Concat(Range(1, 4), Throw(xErrTest)).Pipe(addLatency2)),
				toString,
			)
		}
		subscribe(
			t, observables[:],
			"[A 1]", "[B 2]", xComplete,
			"[A 1]", "[B 2]", "[C 3]", xComplete,
			"[A 1]", "[B 2]", "[C 3]", xErrTest,
		)
	}
}
