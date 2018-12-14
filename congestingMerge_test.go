package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_CongestingMergeAll(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(operators.CongestingMergeAll()),
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(CongestingMergeConfigure{ProjectToObservable, 1}.MakeFunc()),
		},
		"E", "C", "A", "F", "D", "B", xComplete,
		"A", "B", "C", "D", "E", "F", xComplete,
	)
}
