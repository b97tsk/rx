package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_DebounceTime(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(
				addLatencyToValue(1, 2),
				operators.DebounceTime(step(3)),
			),
			Just("A", "B", "C").Pipe(
				addLatencyToValue(1, 3),
				operators.DebounceTime(step(2)),
			),
		},
		"C", Complete,
		"A", "B", "C", Complete,
	)
}
