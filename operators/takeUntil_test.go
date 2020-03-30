package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_TakeUntil(t *testing.T) {
	addLatency := addLatencyToValue(0, 2)
	delay := delaySubscription(3)
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Just(42))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Empty())),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Throw(errTest))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Never())),
		},
		[][]interface{}{
			{Complete},
			{Complete},
			{errTest},
			{"A", "B", "C", Complete},
		},
	)
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Just(42).Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Empty().Pipe(delay))),
			Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(Throw(errTest).Pipe(delay))),
		},
		[][]interface{}{
			{"A", "B", Complete},
			{"A", "B", Complete},
			{"A", "B", errTest},
		},
	)
}
