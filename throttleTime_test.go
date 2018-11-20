package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ThrottleTime(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 2),
				operators.ThrottleTime(step(3)),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 4),
				ThrottleTimeOperator{
					Duration: step(9),
					Leading:  false,
					Trailing: true,
				}.MakeFunc(),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 4),
				ThrottleTimeOperator{
					Duration: step(9),
					Leading:  true,
					Trailing: true,
				}.MakeFunc(),
			),
		},
		"A", "C", "E", "G", xComplete,
		"C", "E", "G", xComplete,
		"A", "C", "E", "G", xComplete,
	)
}
