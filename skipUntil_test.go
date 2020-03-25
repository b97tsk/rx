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
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Throw(errTest))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Never())),
		},
		"A", "B", "C", Complete,
		Complete,
		errTest,
		Complete,
	)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Just(42).Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Empty().Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Throw(errTest).Pipe(delay))),
		},
		"C", Complete,
		Complete,
		errTest,
	)
}
