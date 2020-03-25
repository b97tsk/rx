package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Throttle(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 2),
				operators.Throttle(func(interface{}) Observable {
					return Interval(step(3))
				}),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 2),
				operators.Throttle(func(interface{}) Observable {
					return Throw(errTest)
				}),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 4),
				ThrottleConfigure{
					DurationSelector: func(interface{}) Observable {
						return Interval(step(9))
					},
					Leading:  false,
					Trailing: true,
				}.Use(),
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(0, 4),
				ThrottleConfigure{
					DurationSelector: func(interface{}) Observable {
						return Interval(step(9))
					},
					Leading:  true,
					Trailing: true,
				}.Use(),
			),
		},
		"A", "C", "E", "G", Complete,
		"A", errTest,
		"C", "E", "G", Complete,
		"A", "C", "E", "G", Complete,
	)
}
