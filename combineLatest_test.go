package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_CombineAll(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(operators.CombineAll(), toString),
		},
		"[A C E]", "[A C F]", "[A D F]", "[B D F]", xComplete,
	)
}
