package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSwitch(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(
				rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
				rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
				rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
			).Pipe(AddLatencyToValues(0, 5), operators.Switch()),
			rx.Just(
				rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
				rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
				rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
				rx.Empty(),
			).Pipe(AddLatencyToValues(0, 5), operators.Switch()),
			rx.Just(
				rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
				rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
				rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
				rx.Throw(ErrTest),
			).Pipe(AddLatencyToValues(0, 5), operators.Switch()),
		},
		[][]interface{}{
			{"A", "B", "C", "E", "F", "I", "J", "K", "L", rx.Completed},
			{"A", "B", "C", "E", "F", "I", "J", "K", rx.Completed},
			{"A", "B", "C", "E", "F", "I", "J", "K", ErrTest},
		},
	)
}
