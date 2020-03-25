package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ConcatAll(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(operators.ConcatAll()),
		},
		"A", "B", "C", "D", "E", "F", Complete,
	)
}
