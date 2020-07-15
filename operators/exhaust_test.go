package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestExhaustAll(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(
				rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
				rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
				rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
			).Pipe(AddLatencyToValues(0, 5), operators.ExhaustAll()),
			rx.Just(
				rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
				rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
				rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
				rx.Throw(ErrTest),
			).Pipe(AddLatencyToValues(0, 5), operators.ExhaustAll()),
			rx.Just(
				rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
				rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
				rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
				rx.Throw(ErrTest),
				rx.Throw(ErrTest),
			).Pipe(AddLatencyToValues(0, 5), operators.ExhaustAll()),
		},
		[][]interface{}{
			{"A", "B", "C", "D", "I", "J", "K", "L", Completed},
			{"A", "B", "C", "D", "I", "J", "K", "L", Completed},
			{"A", "B", "C", "D", "I", "J", "K", "L", ErrTest},
		},
	)
}
