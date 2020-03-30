package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestAudit(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Interval(Step(3))
			}),
		),
		"B", "D", "F", rx.Complete,
	)
}
