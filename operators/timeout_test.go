package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTimeout(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(AddLatencyToValues(1, 1), operators.Timeout(Step(2))),
			rx.Just("A", "B", "C").Pipe(AddLatencyToValues(1, 3), operators.Timeout(Step(2))),
		},
		[][]interface{}{
			{"A", "B", "C", Completed},
			{"A", rx.ErrTimeout},
		},
	)
}
