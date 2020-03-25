package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Exhaust(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just(
				Just("A", "B", "C", "D").Pipe(addLatencyToValue(0, 2)),
				Just("E", "F", "G", "H").Pipe(addLatencyToValue(0, 3)),
				Just("I", "J", "K", "L").Pipe(addLatencyToValue(0, 2)),
			).Pipe(addLatencyToValue(0, 5), operators.Exhaust()),
			Just(
				Just("A", "B", "C", "D").Pipe(addLatencyToValue(0, 2)),
				Just("E", "F", "G", "H").Pipe(addLatencyToValue(0, 3)),
				Just("I", "J", "K", "L").Pipe(addLatencyToValue(0, 2)),
				Throw(errTest),
			).Pipe(addLatencyToValue(0, 5), operators.Exhaust()),
			Just(
				Just("A", "B", "C", "D").Pipe(addLatencyToValue(0, 2)),
				Just("E", "F", "G", "H").Pipe(addLatencyToValue(0, 3)),
				Just("I", "J", "K", "L").Pipe(addLatencyToValue(0, 2)),
				Throw(errTest),
				Throw(errTest),
			).Pipe(addLatencyToValue(0, 5), operators.Exhaust()),
		},
		"A", "B", "C", "D", "I", "J", "K", "L", Complete,
		"A", "B", "C", "D", "I", "J", "K", "L", Complete,
		"A", "B", "C", "D", "I", "J", "K", "L", errTest,
	)
}
