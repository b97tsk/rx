package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_TakeUntil(t *testing.T) {
	addLatency := addLatencyToValue(0, 2)
	delay := delaySubscription(3)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Just(42))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Empty())),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Throw(xErrTest))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Never())),
		},
		xComplete,
		xComplete,
		xErrTest,
		"A", "B", "C", xComplete,
	)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Just(42).Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Empty().Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Throw(xErrTest).Pipe(delay))),
		},
		"A", "B", xComplete,
		"A", "B", xComplete,
		"A", "B", xErrTest,
	)
}
