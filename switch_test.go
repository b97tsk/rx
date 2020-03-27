package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Switch(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just(
				Just("A", "B", "C", "D").Pipe(addLatencyToValue(0, 2)),
				Just("E", "F", "G", "H").Pipe(addLatencyToValue(0, 3)),
				Just("I", "J", "K", "L").Pipe(addLatencyToValue(0, 2)),
			).Pipe(addLatencyToValue(0, 5), operators.Switch()),
			Just(
				Just("A", "B", "C", "D").Pipe(addLatencyToValue(0, 2)),
				Just("E", "F", "G", "H").Pipe(addLatencyToValue(0, 3)),
				Just("I", "J", "K", "L").Pipe(addLatencyToValue(0, 2)),
				Empty(),
			).Pipe(addLatencyToValue(0, 5), operators.Switch()),
			Just(
				Just("A", "B", "C", "D").Pipe(addLatencyToValue(0, 2)),
				Just("E", "F", "G", "H").Pipe(addLatencyToValue(0, 3)),
				Just("I", "J", "K", "L").Pipe(addLatencyToValue(0, 2)),
				Throw(errTest),
			).Pipe(addLatencyToValue(0, 5), operators.Switch()),
		},
		[][]interface{}{
			{"A", "B", "C", "E", "F", "I", "J", "K", "L", Complete},
			{"A", "B", "C", "E", "F", "I", "J", "K", Complete},
			{"A", "B", "C", "E", "F", "I", "J", "K", errTest},
		},
	)
}
