package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SampleTime(t *testing.T) {
	subscribe(
		t,
		Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			addLatencyToValue(1, 2),
			operators.SampleTime(step(4)),
		),
		"B", "D", "F", Complete,
	)
}
