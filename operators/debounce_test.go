package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Debounce(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(
				addLatencyToValue(1, 2),
				operators.Debounce(func(interface{}) Observable {
					return Interval(step(3))
				}),
			),
			Just("A", "B", "C").Pipe(
				addLatencyToValue(1, 3),
				operators.Debounce(func(interface{}) Observable {
					return Interval(step(2))
				}),
			),
		},
		[][]interface{}{
			{"C", Complete},
			{"A", "B", "C", Complete},
		},
	)
}
