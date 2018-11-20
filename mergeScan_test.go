package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_MergeScan(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(operators.MergeScan(
				func(acc, val interface{}) Observable {
					return val.(Observable)
				},
				nil,
			)),
		},
		"E", "C", "A", "F", "D", "B", xComplete,
	)
}
