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
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 2),
				operators.Throttle(func(interface{}) rx.Observable {
					return rx.Timer(Step(3))
				}),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 2),
				operators.Throttle(func(interface{}) rx.Observable {
					return rx.Empty()
				}),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 2),
				operators.ThrottleConfigure{
					DurationSelector: func(interface{}) rx.Observable {
						return rx.Empty().Pipe(DelaySubscription(5))
					},
					Leading:  false,
					Trailing: true,
				}.Make(),
			),
			rx.Throw(ErrTest).Pipe(
				operators.Throttle(func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				}),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 2),
				operators.Throttle(func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				}),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleConfigure{
					DurationSelector: func(interface{}) rx.Observable {
						return rx.Timer(Step(9))
					},
					Leading:  false,
					Trailing: true,
				}.Make(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleConfigure{
					DurationSelector: func(interface{}) rx.Observable {
						return rx.Timer(Step(9))
					},
					Leading:  true,
					Trailing: true,
				}.Make(),
			),
		},
		[][]interface{}{
			{"A", "C", "E", Completed},
			{"A", "B", "C", "D", "E", Completed},
			{Completed},
			{ErrTest},
			{"A", ErrTest},
			{"C", "E", Completed},
			{"A", "C", "E", Completed},
		},
	)
}

func TestThrottleTime(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 2),
				operators.ThrottleTime(Step(3)),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleTimeConfigure{
					Duration: Step(9),
					Leading:  false,
					Trailing: true,
				}.Make(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(0, 4),
				operators.ThrottleTimeConfigure{
					Duration: Step(9),
					Leading:  true,
					Trailing: true,
				}.Make(),
			),
		},
		[][]interface{}{
			{"A", "C", "E", Completed},
			{"C", "E", Completed},
			{"A", "C", "E", Completed},
		},
	)
}
