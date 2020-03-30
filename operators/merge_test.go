package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_MergeAll(t *testing.T) {
	subscribe(
		t,
		Just(
			Just("A", "B").Pipe(addLatencyToValue(3, 5)),
			Just("C", "D").Pipe(addLatencyToValue(2, 4)),
			Just("E", "F").Pipe(addLatencyToValue(1, 3)),
		).Pipe(operators.MergeAll()),
		"E", "C", "A", "F", "D", "B", Complete,
	)
}
