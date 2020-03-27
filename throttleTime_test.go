package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ThrottleTime(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 2),
				operators.ThrottleTime(step(3)),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 4),
				ThrottleTimeConfigure{
					Duration: step(9),
					Leading:  false,
					Trailing: true,
				}.Use(),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 4),
				ThrottleTimeConfigure{
					Duration: step(9),
					Leading:  true,
					Trailing: true,
				}.Use(),
			),
		},
		[][]interface{}{
			{"A", "C", "E", "G", Complete},
			{"C", "E", "G", Complete},
			{"A", "C", "E", "G", Complete},
		},
	)
}
