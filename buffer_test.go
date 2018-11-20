package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Buffer(t *testing.T) {
	addLatency := addLatencyToValue(1, 2)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatency,
				operators.Buffer(Interval(step(2))),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatency,
				operators.Buffer(Interval(step(4))),
				toString,
			),
		},
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", xComplete,
		"[A B]", "[C D]", "[E F]", xComplete,
	)
}
