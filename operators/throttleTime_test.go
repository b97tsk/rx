package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestThrottleTime(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 2),
				operators.ThrottleTime(Step(3)),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleTimeConfigure{
					Duration: Step(9),
					Leading:  false,
					Trailing: true,
				}.Use(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleTimeConfigure{
					Duration: Step(9),
					Leading:  true,
					Trailing: true,
				}.Use(),
			),
		},
		[][]interface{}{
			{"A", "C", "E", "G", rx.Complete},
			{"C", "E", "G", rx.Complete},
			{"A", "C", "E", "G", rx.Complete},
		},
	)
}
