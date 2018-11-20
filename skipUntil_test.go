package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SkipUntil(t *testing.T) {
	addLatency := addLatencyToValue(0, 2)
	delay := delaySubscription(3)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Just(42))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Empty())),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Throw(xErrTest))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Never())),
		},
		"A", "B", "C", xComplete,
		xComplete,
		xErrTest,
		xComplete,
	)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Just(42).Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Empty().Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Throw(xErrTest).Pipe(delay))),
		},
		"C", xComplete,
		xComplete,
		xErrTest,
	)
}
