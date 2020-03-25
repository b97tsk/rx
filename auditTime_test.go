package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_AuditTime(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.AuditTime(step(3)),
			),
		},
		"B", "D", "F", Complete,
	)
}
