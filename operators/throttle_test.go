package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestThrottle(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 2),
				operators.Throttle(func(interface{}) rx.Observable {
					return rx.Timer(Step(3))
				}),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 2),
				operators.Throttle(func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				}),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleConfigure{
					DurationSelector: func(interface{}) rx.Observable {
						return rx.Timer(Step(9))
					},
					Leading:  false,
					Trailing: true,
				}.Use(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleConfigure{
					DurationSelector: func(interface{}) rx.Observable {
						return rx.Timer(Step(9))
					},
					Leading:  true,
					Trailing: true,
				}.Use(),
			),
		},
		[][]interface{}{
			{"A", "C", "E", "G", rx.Completed},
			{"A", ErrTest},
			{"C", "E", "G", rx.Completed},
			{"A", "C", "E", "G", rx.Completed},
		},
	)
}

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
			{"A", "C", "E", "G", rx.Completed},
			{"C", "E", "G", rx.Completed},
			{"A", "C", "E", "G", rx.Completed},
		},
	)
}
