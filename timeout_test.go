package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Timeout(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(addLatencyToValue(1, 1), operators.Timeout(step(2))),
			Just("A", "B", "C").Pipe(addLatencyToValue(1, 3), operators.Timeout(step(2))),
		},
		"A", "B", "C", xComplete,
		"A", ErrTimeout,
	)
}
