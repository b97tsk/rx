package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestCongestingMergeAll(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(
				rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
				rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
				rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
			).Pipe(operators.CongestingMergeAll()),
			rx.Just(
				rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
				rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
				rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
			).Pipe(operators.CongestingMergeConfigure{rx.ProjectToObservable, 1}.Use()),
		},
		[][]interface{}{
			{"E", "C", "A", "F", "D", "B", rx.Complete},
			{"A", "B", "C", "D", "E", "F", rx.Complete},
		},
	)
}
