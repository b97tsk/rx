package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestDebounce(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 2),
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Interval(Step(3))
				}),
			),
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 3),
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Interval(Step(2))
				}),
			),
		},
		[][]interface{}{
			{"C", rx.Complete},
			{"A", "B", "C", rx.Complete},
		},
	)
}
