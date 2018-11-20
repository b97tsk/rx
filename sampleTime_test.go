package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SampleTime(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.SampleTime(step(4)),
			),
		},
		"B", "D", "F", xComplete,
	)
}
