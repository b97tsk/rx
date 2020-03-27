package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SkipUntil(t *testing.T) {
	addLatency := addLatencyToValue(0, 2)
	delay := delaySubscription(3)
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Just(42))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Empty())),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Throw(errTest))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Never())),
		},
		[][]interface{}{
			{"A", "B", "C", Complete},
			{Complete},
			{errTest},
			{Complete},
		},
	)
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Just(42).Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Empty().Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(Throw(errTest).Pipe(delay))),
		},
		[][]interface{}{
			{"C", Complete},
			{Complete},
			{errTest},
		},
	)
}
