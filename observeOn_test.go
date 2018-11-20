package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ObserveOn(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Merge(
				Just("A", "B").Pipe(operators.ObserveOn(step(1))),
				Just("C", "D").Pipe(operators.ObserveOn(step(2))),
				Just("E", "F").Pipe(operators.ObserveOn(step(3))),
			),
		},
		"A", "B", "C", "D", "E", "F", xComplete,
	)
}
