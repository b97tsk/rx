package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Sample(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Sample(Interval(step(4))),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Sample(
					Concat(
						Interval(step(4)).Pipe(operators.Take(3)),
						Throw(xErrTest),
					),
				),
			),
		},
		"B", "D", "F", xComplete,
		"B", "D", "F", xErrTest,
	)
}
